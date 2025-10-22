package application

import (
	"bless-activity/model"
	"fmt"
	"log/slog"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type fixBugHandler func(e *core.BootstrapEvent) error

func (application *Application) fixBug(e *core.BootstrapEvent) error {
	list := []fixBugHandler{
		application.fixExample,
		//application.rewardReissue,
		//application.retryFailedPoints,
		//application.articleScoreAndReward,
	}

	for _, handler := range list {
		if err := handler(e); err != nil {
			return err
		}
	}

	return nil
}

func (application *Application) fixExample(*core.BootstrapEvent) error {
	return nil
}

// 奖励补发
func (application *Application) rewardReissue(event *core.BootstrapEvent) error {
	logger := event.App.Logger().With("fix", "rewardReissue")

	// 1. 预加载所有奖励数据到缓存
	rewardCache := make(map[string]*model.Reward)
	var allRewards []*model.Reward
	if err := event.App.RecordQuery(model.DbNameRewards).All(&allRewards); err != nil {
		logger.Error("预加载奖励数据失败", slog.Any("err", err))
		return err
	}
	for _, r := range allRewards {
		rewardCache[r.Id] = r
	}

	// 2. 预加载所有奖项数据到缓存
	awardCache := make(map[string]*model.Awards)
	var allAwards []*model.Awards
	if err := event.App.RecordQuery(model.DbNameAwards).All(&allAwards); err != nil {
		logger.Error("预加载奖项数据失败", slog.Any("err", err))
		return err
	}
	for _, a := range allAwards {
		awardCache[a.Id] = a
	}

	// 3. 预加载所有用户数据到缓存
	userCache := make(map[string]*model.User)
	var allUsers []*model.User
	if err := event.App.RecordQuery(model.DbNameUsers).All(&allUsers); err != nil {
		logger.Error("预加载用户数据失败", slog.Any("err", err))
		return err
	}
	for _, u := range allUsers {
		userCache[u.Id] = u
	}

	// 4. 查找所有 gotReward = false 的历史记录，按创建时间升序排序（先到先得）
	var histories []*model.Histories
	if err := event.App.RecordQuery(model.DbNameHistories).
		Where(dbx.HashExp{model.HistoriesFieldGotReward: false}).
		AndWhere(dbx.Not(dbx.HashExp{model.HistoriesFieldRewardId: ""})).
		OrderBy(model.HistoriesFieldCreated + " asc").
		All(&histories); err != nil {
		logger.Error("查找历史记录失败", slog.Any("err", err))
		return err
	}

	logger.Info("开始补发奖励", slog.Int("count", len(histories)))

	// 5. 统计每个奖励已经补发的数量（在本次补发过程中动态维护）
	rewardIssuedCount := make(map[string]int)
	for rewardId := range rewardCache {
		count, err := event.App.CountRecords(model.DbNameHistories, dbx.HashExp{
			model.HistoriesFieldRewardId:  rewardId,
			model.HistoriesFieldGotReward: true,
		})
		if err != nil {
			logger.Warn("查询已发放数量失败", slog.String("reward_id", rewardId), slog.Any("err", err))
			rewardIssuedCount[rewardId] = 0
		} else {
			rewardIssuedCount[rewardId] = int(count)
		}
	}

	// 6. 获取 points collection
	pointsCollection, err := event.App.FindCollectionByNameOrId(model.DbNamePoints)
	if err != nil {
		logger.Error("查找points集合失败", slog.Any("err", err))
		return err
	}

	successCount := 0
	skipCount := 0

	// 7. 遍历历史记录进行补发
	for _, history := range histories {
		// 从缓存获取奖励信息
		reward, exists := rewardCache[history.RewardId()]
		if !exists {
			logger.Warn("奖励不存在", slog.String("reward_id", history.RewardId()))
			skipCount++
			continue
		}

		// 检查该奖励是否还有余量
		currentIssued := rewardIssuedCount[reward.Id]
		if currentIssued >= reward.Amount() {
			logger.Debug("奖励已发完",
				slog.String("reward_id", reward.Id),
				slog.String("reward_name", reward.Name()),
				slog.Int("issued", currentIssued),
				slog.Int("amount", reward.Amount()))
			skipCount++
			continue
		}

		// 检查状元级别的特殊规则：必须是 isBest 才能获得
		if history.IsTop() && !history.IsBest() {
			logger.Debug("状元级别奖励需要isBest才能获得",
				slog.String("history_id", history.Id),
				slog.Bool("is_best", history.IsBest()))
			skipCount++
			continue
		}

		// 更新历史记录为已获得奖励
		history.SetGotReward(true)
		if err := event.App.Save(history); err != nil {
			logger.Error("更新历史记录失败", slog.String("history_id", history.Id), slog.Any("err", err))
			skipCount++
			continue
		}

		// 更新已发放计数
		rewardIssuedCount[reward.Id]++

		// 如果有积分奖励，创建积分订单并发放
		if reward.Point() > 0 {
			// 从缓存获取奖项名称
			awardName := ""
			if history.AwardId() != "" {
				if award, exists := awardCache[history.AwardId()]; exists {
					awardName = award.Name()
				}
			}

			// 创建积分订单
			pointsRecord := model.NewPointsFromCollection(pointsCollection)
			pointsRecord.SetUserId(history.UserId())
			pointsRecord.SetHistoryId(history.Id)
			pointsRecord.SetPoint(reward.Point())
			pointsRecord.SetStatus(model.PointStatusPending)
			pointsRecord.SetMemo(fmt.Sprintf("【补发】活动《双节同庆·福签传情》第%d次博饼：%s(%s)",
				history.Times(), awardName, reward.Name()))

			if err := event.App.Save(pointsRecord); err != nil {
				logger.Error("保存积分订单失败", slog.Any("err", err))
				continue
			}

			// 从缓存获取用户信息
			user, exists := userCache[history.UserId()]
			if !exists {
				logger.Error("用户不存在", slog.String("user_id", history.UserId()))
				pointsRecord.SetStatus(model.PointStatusFailed)
				pointsRecord.SetError("用户不存在")
				_ = event.App.Save(pointsRecord)
				continue
			}

			// 发放积分
			memo := fmt.Sprintf("%s 交易单号：%s", pointsRecord.Memo(), pointsRecord.Id)
			var distributeErr error
			if !event.App.IsDev() {
				distributeErr = application.fishPiService.Distribute(user.Name(), reward.Point(), memo)
			}

			if distributeErr != nil {
				logger.Error("发放积分失败",
					slog.String("user", user.Name()),
					slog.Int("point", reward.Point()),
					slog.Any("err", distributeErr))
				pointsRecord.SetStatus(model.PointStatusFailed)
				pointsRecord.SetError(distributeErr.Error())
			} else {
				pointsRecord.SetStatus(model.PointStatusSuccess)
				successCount++
				logger.Info("补发积分成功",
					slog.String("user", user.Name()),
					slog.String("reward", reward.Name()),
					slog.Int("point", reward.Point()),
					slog.Int("times", history.Times()))
			}

			if err := event.App.Save(pointsRecord); err != nil {
				logger.Error("更新积分订单状态失败", slog.Any("err", err))
			}
		} else {
			// 没有积分奖励，也算作成功处理
			successCount++
		}
	}

	logger.Info("奖励补发完成",
		slog.Int("total", len(histories)),
		slog.Int("success", successCount),
		slog.Int("skip", skipCount))
	return nil
}

// 重新发放失败的积分订单
func (application *Application) retryFailedPoints(event *core.BootstrapEvent) error {
	logger := event.App.Logger().With("fix", "retryFailedPoints")

	// 预加载所有用户数据到缓存
	userCache := make(map[string]*model.User)
	var allUsers []*model.User
	if err := event.App.RecordQuery(model.DbNameUsers).All(&allUsers); err != nil {
		logger.Error("预加载用户数据失败", slog.Any("err", err))
		return err
	}
	for _, u := range allUsers {
		userCache[u.Id] = u
	}

	// 查找所有失败的积分订单
	var failedPoints []*model.Points
	if err := event.App.RecordQuery(model.DbNamePoints).
		Where(dbx.HashExp{model.PointsFieldStatus: model.PointStatusFailed}).
		OrderBy(model.PointsFieldCreated + " asc").
		All(&failedPoints); err != nil {
		logger.Error("查找失败的积分订单失败", slog.Any("err", err))
		return err
	}

	logger.Info("开始重新发放失败的积分订单", slog.Int("count", len(failedPoints)))

	successCount := 0
	failCount := 0

	for i, pointsRecord := range failedPoints {
		// 从缓存获取用户信息
		user, exists := userCache[pointsRecord.UserId()]
		if !exists {
			logger.Warn("用户不存在，跳过",
				slog.String("user_id", pointsRecord.UserId()),
				slog.String("points_id", pointsRecord.Id))
			failCount++
			continue
		}

		// 更新订单状态为待发放
		pointsRecord.SetStatus(model.PointStatusPending)
		pointsRecord.SetError("") // 清空之前的错误信息

		// 重新尝试发放积分
		memo := fmt.Sprintf("%s 交易单号：%s", pointsRecord.Memo(), pointsRecord.Id)
		var distributeErr error
		if !event.App.IsDev() {
			distributeErr = application.fishPiService.Distribute(user.Name(), pointsRecord.Point(), memo)
		}

		if distributeErr != nil {
			// 发放失败
			logger.Error("重新发放积分失败",
				slog.String("user", user.Name()),
				slog.Int("point", pointsRecord.Point()),
				slog.String("points_id", pointsRecord.Id),
				slog.Any("err", distributeErr))
			pointsRecord.SetStatus(model.PointStatusFailed)
			pointsRecord.SetError(distributeErr.Error())
			failCount++
		} else {
			// 发放成功
			pointsRecord.SetStatus(model.PointStatusSuccess)
			successCount++
			logger.Info("重新发放积分成功",
				slog.String("user", user.Name()),
				slog.Int("point", pointsRecord.Point()),
				slog.String("points_id", pointsRecord.Id))
		}

		// 保存订单状态
		if err := event.App.Save(pointsRecord); err != nil {
			logger.Error("更新积分订单状态失败",
				slog.String("points_id", pointsRecord.Id),
				slog.Any("err", err))
		}

		// 每处理一条记录后延迟一段时间，避免请求过快
		// 每10条延迟1秒，避免API限流
		if (i+1)%10 == 0 {
			logger.Info("已处理10条记录，等待1秒",
				slog.Int("processed", i+1),
				slog.Int("total", len(failedPoints)))
			time.Sleep(1 * time.Second)
		} else {
			// 其他情况延迟100ms
			time.Sleep(200 * time.Millisecond)
		}
	}

	logger.Info("失败积分订单重新发放完成",
		slog.Int("total", len(failedPoints)),
		slog.Int("success", successCount),
		slog.Int("failed", failCount))

	return nil
}

// 文章评分和奖励发放
func (application *Application) articleScoreAndReward(event *core.BootstrapEvent) error {
	logger := event.App.Logger().With("fix", "articleScoreAndReward")

	// 1. 获取所有文章
	var articles []*model.Article
	if err := event.App.RecordQuery(model.DbNameArticles).Where(dbx.Not(dbx.HashExp{model.CommonFieldId: "4tmbyzsuhbd66yr"})).All(&articles); err != nil {
		logger.Error("查询文章失败", slog.Any("err", err))
		return err
	}

	logger.Info("开始计算文章评分", slog.Int("total", len(articles)))

	// 2. 计算并更新每篇文章的评分
	type ArticleScore struct {
		Article *model.Article
		Score   float64
	}

	articleScores := make([]ArticleScore, 0, len(articles))

	for _, article := range articles {
		score := article.CalculateScore()
		article.SetScore(score)

		if err := event.App.Save(article); err != nil {
			logger.Error("更新文章评分失败",
				slog.String("article_id", article.Id),
				slog.String("title", article.Title()),
				slog.Any("err", err))
			continue
		}

		articleScores = append(articleScores, ArticleScore{
			Article: article,
			Score:   score,
		})

		logger.Debug("文章评分计算完成",
			slog.String("title", article.Title()),
			slog.Float64("score", score),
			slog.Int("viewCount", article.ViewCount()),
			slog.Int("goodCnt", article.GoodCnt()),
			slog.Int("collectCnt", article.CollectCnt()),
			slog.Int("commentCount", article.CommentCount()),
			slog.Int("thankCnt", article.ThankCnt()))
	}

	// 3. 按评分排序（从高到低）
	for i := 0; i < len(articleScores); i++ {
		for j := i + 1; j < len(articleScores); j++ {
			if articleScores[j].Score > articleScores[i].Score {
				articleScores[i], articleScores[j] = articleScores[j], articleScores[i]
			}
		}
	}

	logger.Info("文章评分排序完成", slog.Int("count", len(articleScores)))

	// 4. 预加载用户缓存
	userCache := make(map[string]*model.User)
	var allUsers []*model.User
	if err := event.App.RecordQuery(model.DbNameUsers).All(&allUsers); err != nil {
		logger.Error("预加载用户数据失败", slog.Any("err", err))
		return err
	}
	for _, u := range allUsers {
		userCache[u.Id] = u
	}

	// 5. 获取 points collection
	pointsCollection, err := event.App.FindCollectionByNameOrId(model.DbNamePoints)
	if err != nil {
		logger.Error("查找points集合失败", slog.Any("err", err))
		return err
	}

	// 6. 根据排名发放积分
	successCount := 0
	failCount := 0

	for rank, as := range articleScores {
		ranking := rank + 1
		var points int
		var rankName string

		// 根据排名确定积分
		if ranking == 1 {
			points = 1024
			rankName = "第一名"
		} else if ranking >= 2 && ranking <= 3 {
			points = 512
			rankName = fmt.Sprintf("第%d名", ranking)
		} else if ranking >= 4 && ranking <= 10 {
			points = 256
			rankName = fmt.Sprintf("第%d名", ranking)
		} else {
			points = 128
			rankName = "参与奖"
		}

		// 获取文章作者
		user, exists := userCache[as.Article.UserId()]
		if !exists {
			logger.Warn("文章作者不存在，跳过",
				slog.String("article_id", as.Article.Id),
				slog.String("user_id", as.Article.UserId()))
			failCount++
			continue
		}

		// 创建积分订单
		pointsRecord := model.NewPointsFromCollection(pointsCollection)
		pointsRecord.SetUserId(user.Id)
		pointsRecord.SetPoint(points)
		pointsRecord.SetStatus(model.PointStatusPending)
		pointsRecord.SetMemo(fmt.Sprintf("活动《双节同庆·福签传情》文章评分奖励：%s（评分：%.2f）",
			rankName, as.Score))

		if err := event.App.Save(pointsRecord); err != nil {
			logger.Error("保存积分订单失败",
				slog.String("user", user.Name()),
				slog.Any("err", err))
			failCount++
			continue
		}

		// 发放积分
		memo := fmt.Sprintf("%s 交易单号：%s", pointsRecord.Memo(), pointsRecord.Id)
		var distributeErr error
		if !event.App.IsDev() {
			distributeErr = application.fishPiService.Distribute(user.Name(), points, memo)
		}

		if distributeErr != nil {
			logger.Error("发放积分失败",
				slog.String("user", user.Name()),
				slog.Int("point", points),
				slog.Any("err", distributeErr))
			pointsRecord.SetStatus(model.PointStatusFailed)
			pointsRecord.SetError(distributeErr.Error())
			failCount++
		} else {
			pointsRecord.SetStatus(model.PointStatusSuccess)
			successCount++
			logger.Info("发放积分成功",
				slog.Int("ranking", ranking),
				slog.String("user", user.Name()),
				slog.String("article", as.Article.Title()),
				slog.Float64("score", as.Score),
				slog.Int("point", points))
		}

		if err := event.App.Save(pointsRecord); err != nil {
			logger.Error("更新积分订单状态失败", slog.Any("err", err))
		}

		// 延迟，避免请求过快
		if ranking%10 == 0 {
			logger.Info("已处理10条记录，等待1秒",
				slog.Int("processed", ranking),
				slog.Int("total", len(articleScores)))
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(200 * time.Millisecond)
		}
	}

	logger.Info("文章评分奖励发放完成",
		slog.Int("total", len(articleScores)),
		slog.Int("success", successCount),
		slog.Int("failed", failCount))

	return nil
}
