package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type VoteController struct {
	event *core.ServeEvent
	app   core.App

	logger *slog.Logger
}

func NewVoteController(event *core.ServeEvent) *VoteController {
	logger := event.App.Logger().With(
		slog.String("controller", "vote"),
	)

	controller := &VoteController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *VoteController) registerRoutes() {
	group := controller.event.Router.Group("/vote")
	group.POST("", controller.CreateVote).BindFunc(controller.CheckLogin)
	group.DELETE("/{id}", controller.DeleteVote).BindFunc(controller.CheckLogin)
	group.GET("/my", controller.GetMyVotes).BindFunc(controller.CheckLogin)
	group.GET("/rank", controller.GetVoteRank)
	group.GET("/statistics", controller.GetStatistics)
}

func (controller *VoteController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *VoteController) CheckLogin(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	if event.HasSuperuserAuth() {
		return event.ForbiddenError("请登录普通用户账号", nil)
	}
	return event.Next()
}

// CreateVote 创建投票（赠送福签）
func (controller *VoteController) CreateVote(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("create_vote")

	user := model.NewUser(event.Auth)

	// 解析请求体
	data := struct {
		ArticleId string `json:"article_id"`
		VoteType  string `json:"vote_type"`
	}{}

	if err := event.BindBody(&data); err != nil {
		return event.BadRequestError("请求参数错误", err)
	}

	// 验证投票类型
	if data.VoteType != model.VoteTypeCareer && data.VoteType != model.VoteTypeRomance && data.VoteType != model.VoteTypeWealth {
		return event.BadRequestError("无效的福签类型", nil)
	}

	// 查找目标文章
	article := new(model.Article)
	if err := controller.app.RecordQuery(model.DbNameArticles).
		Where(dbx.HashExp{model.CommonFieldId: data.ArticleId}).
		One(article); err != nil {
		logger.Error("查找文章失败", slog.Any("err", err))
		return event.NotFoundError("文章不存在", err)
	}

	// 不能给自己投票
	if article.UserId() == user.Id {
		return event.BadRequestError("不能给自己的文章赠送福签", nil)
	}

	// 检查是否已经给该用户投过这种类型的票
	existingVote := new(model.Vote)
	err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{
			model.VotesFieldFromUserId: user.Id,
			model.VotesFieldVoteType:   data.VoteType,
		}).
		One(existingVote)

	if err == nil {
		// 已存在该类型的投票
		return event.BadRequestError("您已经赠送过这种福签了", nil)
	}

	// 创建投票记录
	votesCollection, err := controller.app.FindCollectionByNameOrId(model.DbNameVotes)
	if err != nil {
		logger.Error("查找votes集合失败", slog.Any("err", err))
		return event.InternalServerError("查找votes集合失败", err)
	}

	vote := model.NewVoteFromCollection(votesCollection)
	vote.SetFromUserId(user.Id)
	vote.SetToUserId(article.UserId())
	vote.SetArticleId(article.Id)
	vote.SetVoteType(data.VoteType)

	if err := controller.app.Save(vote); err != nil {
		logger.Error("保存投票记录失败", slog.Any("err", err))
		return event.InternalServerError("保存投票记录失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "福签赠送成功",
		"vote_id": vote.Id,
	})
}

// DeleteVote 撤销投票
func (controller *VoteController) DeleteVote(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("delete_vote")

	user := model.NewUser(event.Auth)
	voteId := event.Request.PathValue("id")

	if voteId == "" {
		return event.BadRequestError("投票ID不能为空", nil)
	}

	// 查找投票记录
	vote := new(model.Vote)
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.CommonFieldId: voteId}).
		One(vote); err != nil {
		logger.Error("查找投票记录失败", slog.Any("err", err))
		return event.NotFoundError("投票记录不存在", err)
	}

	// 验证是否为本人的投票
	if vote.FromUserId() != user.Id {
		return event.ForbiddenError("不能撤销他人的投票", nil)
	}

	// 删除投票记录
	if err := controller.app.Delete(vote); err != nil {
		logger.Error("删除投票记录失败", slog.Any("err", err))
		return event.InternalServerError("删除投票记录失败", err)
	}

	return event.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "投票撤销成功",
	})
}

// GetMyVotes 获取我的投票记录
func (controller *VoteController) GetMyVotes(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_my_votes")

	user := model.NewUser(event.Auth)

	// 查询我的投票记录
	votes := []*model.Vote{}
	if err := controller.app.RecordQuery(model.DbNameVotes).
		Where(dbx.HashExp{model.VotesFieldFromUserId: user.Id}).
		OrderBy(model.VotesFieldCreated + " desc").
		All(&votes); err != nil {
		logger.Error("查找投票记录失败", slog.Any("err", err))
		return event.InternalServerError("查找投票记录失败", err)
	}

	// 构建响应数据
	result := make([]map[string]any, 0, len(votes))
	for _, vote := range votes {
		// 查找目标用户信息
		toUser := new(model.User)
		if err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: vote.ToUserId()}).
			One(toUser); err != nil {
			logger.Warn("查找目标用户失败", slog.Any("err", err))
			continue
		}

		// 查找文章信息
		article := new(model.Article)
		if err := controller.app.RecordQuery(model.DbNameArticles).
			Where(dbx.HashExp{model.CommonFieldId: vote.ArticleId()}).
			One(article); err != nil {
			logger.Warn("查找文章失败", slog.Any("err", err))
			continue
		}

		result = append(result, map[string]any{
			"id":             vote.Id,
			"vote_type":      vote.VoteType(),
			"to_user_name":   toUser.Name(),
			"to_user_nick":   toUser.Nickname(),
			"to_user_avatar": toUser.Avatar(),
			"article_id":     article.Id,
			"article_title":  article.Title(),
			"article_oid":    article.OId(),
			"created":        vote.Created(),
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
	})
}

// GetVoteRank 获取投票排行榜
func (controller *VoteController) GetVoteRank(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_vote_rank")

	// 查询所有投票记录并按接收者分组统计
	votes := []*model.Vote{}
	if err := controller.app.RecordQuery(model.DbNameVotes).All(&votes); err != nil {
		logger.Error("查找投票记录失败", slog.Any("err", err))
		return event.InternalServerError("查找投票记录失败", err)
	}

	// 统计每个用户收到的各类型福签数量
	userStats := make(map[string]map[string]int)
	for _, vote := range votes {
		if _, exists := userStats[vote.ToUserId()]; !exists {
			userStats[vote.ToUserId()] = map[string]int{
				model.VoteTypeCareer:  0,
				model.VoteTypeRomance: 0,
				model.VoteTypeWealth:  0,
				"total":               0,
			}
		}
		userStats[vote.ToUserId()][vote.VoteType()]++
		userStats[vote.ToUserId()]["total"]++
	}

	// 构建排行榜数据
	result := make([]map[string]any, 0)
	for userId, stats := range userStats {
		// 查找用户信息
		user := new(model.User)
		if err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: userId}).
			One(user); err != nil {
			logger.Warn("查找用户失败", slog.Any("err", err))
			continue
		}

		// 查找用户的文章
		article := new(model.Article)
		if err := controller.app.RecordQuery(model.DbNameArticles).
			Where(dbx.HashExp{model.ArticlesFieldUserId: userId}).
			OrderBy(model.ArticlesFieldCreatedAt + " desc").
			One(article); err != nil {
			logger.Warn("查找文章失败", slog.Any("err", err))
			continue
		}

		result = append(result, map[string]any{
			"user_id":       userId,
			"user_name":     user.Name(),
			"user_nickname": user.Nickname(),
			"user_avatar":   user.Avatar(),
			"article_id":    article.Id,
			"article_title": article.Title(),
			"article_oid":   article.OId(),
			"career_count":  stats[model.VoteTypeCareer],
			"romance_count": stats[model.VoteTypeRomance],
			"wealth_count":  stats[model.VoteTypeWealth],
			"total_count":   stats["total"],
		})
	}

	return event.JSON(http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
	})
}

// GetStatistics 获取投票统计信息（用于显示当前用户的投票状态）
func (controller *VoteController) GetStatistics(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_statistics")

	// 如果已登录，返回用户的投票状态
	if event.Auth != nil && !event.HasSuperuserAuth() {
		user := model.NewUser(event.Auth)

		// 查询用户已投的票
		votes := []*model.Vote{}
		if err := controller.app.RecordQuery(model.DbNameVotes).
			Where(dbx.HashExp{model.VotesFieldFromUserId: user.Id}).
			All(&votes); err != nil {
			logger.Error("查找投票记录失败", slog.Any("err", err))
			return event.InternalServerError("查找投票记录失败", err)
		}

		votedTypes := make(map[string]bool)
		for _, vote := range votes {
			votedTypes[vote.VoteType()] = true
		}

		return event.JSON(http.StatusOK, map[string]any{
			"career_voted":  votedTypes[model.VoteTypeCareer],
			"romance_voted": votedTypes[model.VoteTypeRomance],
			"wealth_voted":  votedTypes[model.VoteTypeWealth],
			"remaining":     3 - len(votedTypes),
		})
	}

	// 未登录时返回空状态
	return event.JSON(http.StatusOK, map[string]any{
		"career_voted":  false,
		"romance_voted": false,
		"wealth_voted":  false,
		"remaining":     3,
	})
}
