package controller

import (
	"bless-activity/model"
	"bless-activity/service/fishpi"
	"bless-activity/service/mooncakeGambling"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type MooncakeController struct {
	event *core.ServeEvent
	app   core.App

	logger        *slog.Logger
	game          *mooncakeGambling.MooncakeGame
	fishpiService *fishpi.Service
}

func NewMooncakeController(event *core.ServeEvent, fishpiService *fishpi.Service) *MooncakeController {
	logger := event.App.Logger().With(
		slog.String("controller", "mooncake"),
	)

	controller := &MooncakeController{
		event:         event,
		app:           event.App,
		logger:        logger,
		game:          mooncakeGambling.NewMooncakeGame(),
		fishpiService: fishpiService,
	}

	controller.registerRoutes()

	return controller
}

func (controller *MooncakeController) registerRoutes() {
	group := controller.event.Router.Group("/mooncake")
	group.POST("/gambling", controller.Gambling).BindFunc(controller.CheckLogin)
	group.GET("/history", controller.GetHistory).BindFunc(controller.CheckLogin)
}

func (controller *MooncakeController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *MooncakeController) CheckLogin(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	if event.HasSuperuserAuth() {
		return event.ForbiddenError("请登录普通用户账号", nil)
	}
	return event.Next()
}

// Gambling 博饼接口
func (controller *MooncakeController) Gambling(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("gambling")

	user := model.NewUser(event.Auth)

	// 查找用户最新文章
	article := new(model.Article)
	if err := controller.app.RecordQuery(model.DbNameArticles).
		Where(dbx.HashExp{model.ArticlesFieldUserId: user.Id}).
		OrderBy(model.ArticlesFieldCreatedAt + " desc").
		One(article); err != nil {
		logger.Error("查找最新文章失败", slog.Any("err", err))
		return event.InternalServerError("查找最新文章失败", err)
	}

	// 查询用户已抽奖次数
	drawTimes, err := controller.app.CountRecords(model.DbNameHistories, dbx.HashExp{
		model.HistoriesFieldUserId: user.Id,
	})
	if err != nil {
		logger.Error("查找抽奖次数失败", slog.Any("err", err))
		return event.InternalServerError("查找抽奖次数失败", err)
	}

	// 计算剩余次数
	totalTimes := model.DefaultMooncakeGamblingTimes + article.ThankCnt()
	// 限制最大次数为20次
	if totalTimes > model.MaxMooncakeGamblingTimes {
		totalTimes = model.MaxMooncakeGamblingTimes
	}
	restTimes := totalTimes - int(drawTimes)

	if restTimes <= 0 {
		return event.BadRequestError("博饼次数已用完", nil)
	}

	// 进行博饼
	result := controller.game.Play()

	// 根据 PrizeLevel 查找对应的 awards（可能有多个同等级的奖项）
	selectedAward := new(model.Awards)
	if err = controller.app.RecordQuery(model.DbNameAwards).
		Where(dbx.HashExp{model.AwardsFieldLevel: int(result.PrizeLevel)}).
		One(selectedAward); err != nil {
		logger.Error("查找奖项失败", slog.Any("err", err))
		return event.InternalServerError("查找奖项失败", err)
	}

	// 根据选中的奖项查找对应的 reward
	reward := new(model.Reward)
	if err = controller.app.RecordQuery(model.DbNameRewards).
		Where(dbx.HashExp{model.CommonFieldId: selectedAward.RewardId()}).
		One(reward); err != nil {
		logger.Error("查找奖励记录失败", slog.Any("err", err))
		return event.InternalServerError("查找奖励记录失败", err)
	}

	// 创建历史记录
	historiesCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameHistories)
	if err != nil {
		logger.Error("查找histories集合失败", slog.Any("err", err))
		return event.InternalServerError("查找histories集合失败", err)
	}

	history := model.NewHistoriesFromCollection(historiesCollection)
	history.SetUserId(user.Id)
	history.SetTimes(int(drawTimes) + 1)
	if reward != nil {
		history.SetRewardId(reward.Id)
	}
	if selectedAward != nil {
		history.SetAwardId(selectedAward.Id)
	}
	history.SetIsTop(result.PrizeLevel.IsTop())
	history.SetDetails(result.Dices)

	// 如果是 top 等级，处理 isBest 比较逻辑：
	// - 查找用户最新的一条 isBest 记录（按时间倒序）
	// - 若不存在，则将当前记录标记为 isBest
	// - 若存在，则使用 details() 数据构造 GameResult 并通过 CompareGameResult 比较，较大的为 isBest
	if result.PrizeLevel.IsTop() {
		prevBest := new(model.Histories)
		if err := controller.app.RecordQuery(model.DbNameHistories).
			Where(dbx.HashExp{
				model.HistoriesFieldUserId: user.Id,
				model.HistoriesFieldIsBest: true,
			}).
			OrderBy(model.HistoriesFieldCreated + " desc").
			Limit(1).
			One(prevBest); err != nil {
			// 没有找到之前的 isBest（或查询失败）——把当前标记为 isBest
			history.SetIsBest(true)
		} else {
			// 找到之前的最佳记录，按规则比较两次结果
			prevResult := controller.game.PlayWithDices(prevBest.Details())
			compare := mooncakeGambling.CompareGameResult(result, prevResult)
			if compare > 0 {
				// 当前更好：取消之前的 isBest，并将当前设为 isBest
				prevBest.SetIsBest(false)
				if err := controller.app.Save(prevBest); err != nil {
					logger.Error("取消之前 isBest 记录失败", slog.Any("err", err))
					// 不中断流程
				}
				history.SetIsBest(true)
			} else {
				// 仍不如已有的最佳，当前不是 isBest
				history.SetIsBest(false)
			}
		}
	} else {
		history.SetIsBest(false)
	}

	// 决定是否实际获得奖励（gotReward）
	got := false
	if reward != nil {
		// 查询已经发放的数量（gotReward == true）
		issuedCount, cntErr := controller.app.CountRecords(model.DbNameHistories, dbx.HashExp{
			model.HistoriesFieldRewardId:  reward.Id,
			model.HistoriesFieldGotReward: true,
		})
		if cntErr != nil {
			// 如果查询失败，不中断流程；记录警告，默认不发放
			logger.Warn("查询已发放奖励数量失败，暂不发放奖励", slog.Any("err", cntErr))
			issuedCount = 0
		}

		// 特殊规则：level >= 状元四点红（PrizeLevelZSiDianHong）仅当 isBest 为 true 时可获得
		if result.PrizeLevel >= mooncakeGambling.PrizeLevelZSiDianHong {
			if history.IsBest() && int(issuedCount) < reward.Amount() {
				got = true
			} else {
				got = false
			}
		} else {
			// 其他奖励：先到先得（只要未超出 reward.Amount 就发放）
			if int(issuedCount) < reward.Amount() {
				got = true
			} else {
				got = false
			}
		}
	}

	history.SetGotReward(got)

	if err := controller.app.Save(history); err != nil {
		logger.Error("保存历史记录失败", slog.Any("err", err))
		return event.InternalServerError("保存历史记录失败", err)
	}

	// 发放积分奖励（仅当获得奖励且不是状元级别）
	if got && reward.Point() > 0 && !result.PrizeLevel.IsTop() {
		// 创建积分订单记录，状态为待发放
		pointsCollection, err := controller.app.FindCollectionByNameOrId(model.DbNamePoints)
		if err != nil {
			logger.Error("查找points集合失败", slog.Any("err", err))
			// 不中断流程，只记录错误
		} else {
			pointsRecord := model.NewPointsFromCollection(pointsCollection)
			pointsRecord.SetUserId(user.Id)
			pointsRecord.SetHistoryId(history.Id)
			pointsRecord.SetPoint(reward.Point())
			pointsRecord.SetStatus(model.PointStatusPending)
			pointsRecord.SetMemo(fmt.Sprintf("活动《双节同庆·福签传情》第%d次博饼：%s(%s)", history.Times(), selectedAward.Name(), reward.Name()))

			// 保存订单记录
			if err := controller.app.Save(pointsRecord); err != nil {
				logger.Error("保存积分订单失败", slog.Any("err", err))
			} else {
				// 进行积分发放
				memo := fmt.Sprintf("%s 交易单号：%s", pointsRecord.Memo(), pointsRecord.Id)
				distributeErr := controller.fishpiService.Distribute(user.Name(), reward.Point(), memo)

				if distributeErr != nil {
					// 发放失败，更新订单状态
					logger.Error("发放积分失败", slog.Any("err", distributeErr))
					pointsRecord.SetStatus(model.PointStatusFailed)
					pointsRecord.SetError(distributeErr.Error())
					if err := controller.app.Save(pointsRecord); err != nil {
						logger.Error("更新积分订单状态失败", slog.Any("err", err))
					}
				} else {
					// 发放成功，更新订单状态
					pointsRecord.SetStatus(model.PointStatusSuccess)
					if err := controller.app.Save(pointsRecord); err != nil {
						logger.Error("更新积分订单状态失败", slog.Any("err", err))
					}
				}
			}
		}
	}

	return event.JSON(http.StatusOK, map[string]any{
		"dices":       result.Dices,
		"prize_level": int(result.PrizeLevel),
		"prize_name":  result.PrizeName,
		"rest_times":  restTimes - 1,
		"times":       history.Times(),
		"reward_id": func() string {
			if reward != nil {
				return reward.Id
			}
			return ""
		}(),
		"reward_name": func() string {
			if reward != nil {
				return reward.Name()
			}
			return ""
		}(),
		"award_id": func() string {
			if selectedAward != nil {
				return selectedAward.Id
			}
			return ""
		}(),
		"award_name": func() string {
			if selectedAward != nil {
				return selectedAward.Name()
			}
			return ""
		}(),
		"got_reward": history.GotReward(),
	})
}

// GetHistory 获取博饼历史记录
func (controller *MooncakeController) GetHistory(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_history")

	user := model.NewUser(event.Auth)

	// 查询用户的历史记录（按时间倒序，最多20条）
	histories := []*model.Histories{}
	if err := controller.app.RecordQuery(model.DbNameHistories).
		Where(dbx.HashExp{model.HistoriesFieldUserId: user.Id}).
		OrderBy(model.HistoriesFieldCreated + " desc").
		Limit(20).
		All(&histories); err != nil {
		logger.Error("查找历史记录失败", slog.Any("err", err))
		return event.InternalServerError("查找历史记录失败", err)
	}

	// 转换为响应格式
	result := make([]map[string]any, 0, len(histories))
	for _, h := range histories {
		// 查找对应的reward
		reward := new(model.Reward)
		if err := controller.app.RecordQuery(model.DbNameRewards).
			Where(dbx.HashExp{model.CommonFieldId: h.RewardId()}).
			One(reward); err != nil {
			logger.Warn("查找reward失败", slog.Any("err", err), slog.String("reward_id", h.RewardId()))
			continue
		}

		item := map[string]any{
			"id":          h.Id,
			"times":       h.Times(),
			"dices":       h.Details(),
			"prize_level": reward.Level(),
			"prize_name":  reward.Name(),
			"is_top":      h.IsTop(),
			"is_best":     h.IsBest(),
			"got_reward":  h.GotReward(),
			"created":     h.Created(),
		}

		// 如果有award信息，添加进去
		if h.AwardId() != "" {
			award := new(model.Awards)
			if err := controller.app.RecordQuery(model.DbNameAwards).
				Where(dbx.HashExp{model.CommonFieldId: h.AwardId()}).
				One(award); err == nil {
				item["award_name"] = award.Name()
				item["award_description"] = award.Description()
			}
		}

		result = append(result, item)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
	})
}
