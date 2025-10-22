package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ActivityController struct {
	event *core.ServeEvent
	app   core.App

	logger *slog.Logger
}

func NewActivityController(event *core.ServeEvent) *ActivityController {
	logger := event.App.Logger().With(
		slog.String("controller", "activity"),
	)

	controller := &ActivityController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *ActivityController) registerRoutes() {
	group := controller.event.Router.Group("/activity")
	group.GET("/result", controller.GetActivityResult)
}

func (controller *ActivityController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

// GetActivityResult 获取活动结果数据
func (controller *ActivityController) GetActivityResult(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_activity_result")

	// 1. 获取博饼状元信息（isBest=true的记录）
	var champions []struct {
		UserId   string `db:"userId" json:"user_id"`
		Times    int    `db:"times" json:"times"`
		Details  string `db:"details" json:"details"`
		Created  string `db:"created" json:"created"`
		Username string `db:"username" json:"username"`
		Nickname string `db:"nickname" json:"nickname"`
		Avatar   string `db:"avatar" json:"avatar"`
		AwardId  string `db:"awardId" json:"award_id"`
	}

	err := controller.app.DB().
		NewQuery(`
			SELECT h.userId, h.times, h.details, h.created, h.awardId,
			       u.name as username, u.nickname, u.avatar
			FROM histories h
			LEFT JOIN users u ON h.userId = u.id
			WHERE h.isBest = true
			ORDER BY h.created DESC
		`).
		All(&champions)

	if err != nil {
		logger.Error("查询状元信息失败", slog.Any("err", err))
		return event.InternalServerError("查询状元信息失败", err)
	}

	// 2. 获取博饼参与统计（非状元的参与人员和次数）
	var participants []struct {
		UserId   string `db:"userId" json:"user_id"`
		Count    int    `db:"count" json:"count"`
		Username string `db:"username" json:"username"`
		Nickname string `db:"nickname" json:"nickname"`
		Avatar   string `db:"avatar" json:"avatar"`
	}

	err = controller.app.DB().
		NewQuery(`
			SELECT h.userId, COUNT(*) as count,
			       u.name as username, u.nickname, u.avatar
			FROM histories h
			LEFT JOIN users u ON h.userId = u.id
			WHERE h.isTop = false
			GROUP BY h.userId
			ORDER BY count DESC
		`).
		All(&participants)

	if err != nil {
		logger.Error("查询参与者信息失败", slog.Any("err", err))
		return event.InternalServerError("查询参与者信息失败", err)
	}

	// 3. 获取三种福签的获得信息
	var voteStats []struct {
		ToUserId string `db:"toUserId" json:"to_user_id"`
		VoteType string `db:"voteType" json:"vote_type"`
		Count    int    `db:"count" json:"count"`
		Username string `db:"username" json:"username"`
		Nickname string `db:"nickname" json:"nickname"`
		Avatar   string `db:"avatar" json:"avatar"`
	}

	err = controller.app.DB().
		NewQuery(`
			SELECT v.toUserId, v.voteType, COUNT(*) as count,
			       u.name as username, u.nickname, u.avatar
			FROM votes v
			LEFT JOIN users u ON v.toUserId = u.id
			GROUP BY v.toUserId, v.voteType
			ORDER BY v.voteType, count DESC
		`).
		All(&voteStats)

	if err != nil {
		logger.Error("查询福签信息失败", slog.Any("err", err))
		return event.InternalServerError("查询福签信息失败", err)
	}

	// 按福签类型分组
	careerVotes := []map[string]interface{}{}
	romanceVotes := []map[string]interface{}{}
	wealthVotes := []map[string]interface{}{}

	for _, stat := range voteStats {
		item := map[string]interface{}{
			"user_id":  stat.ToUserId,
			"username": stat.Username,
			"nickname": stat.Nickname,
			"avatar":   stat.Avatar,
			"count":    stat.Count,
		}
		switch stat.VoteType {
		case model.VoteTypeCareer:
			careerVotes = append(careerVotes, item)
		case model.VoteTypeRomance:
			romanceVotes = append(romanceVotes, item)
		case model.VoteTypeWealth:
			wealthVotes = append(wealthVotes, item)
		}
	}

	// 4. 获取文章排名信息
	articles := []*model.Article{}
	err = controller.app.RecordQuery(model.DbNameArticles).
		OrderBy(model.ArticlesFieldScore + " DESC").
		Limit(100).
		All(&articles)

	if err != nil {
		logger.Error("查询文章排名失败", slog.Any("err", err))
		return event.InternalServerError("查询文章排名失败", err)
	}

	// 获取文章对应的用户信息
	articleRankings := []map[string]interface{}{}
	for rank, article := range articles {
		user := new(model.User)
		err := controller.app.RecordQuery(model.DbNameUsers).
			Where(dbx.HashExp{model.CommonFieldId: article.UserId()}).
			One(user)

		if err != nil {
			logger.Warn("查询用户信息失败", slog.String("userId", article.UserId()), slog.Any("err", err))
			continue
		}

		// 统计该文章收到的福签数量
		careerCount, _ := controller.app.CountRecords(model.DbNameVotes, dbx.HashExp{
			model.VotesFieldArticleId: article.Id,
			model.VotesFieldVoteType:  model.VoteTypeCareer,
		})
		romanceCount, _ := controller.app.CountRecords(model.DbNameVotes, dbx.HashExp{
			model.VotesFieldArticleId: article.Id,
			model.VotesFieldVoteType:  model.VoteTypeRomance,
		})
		wealthCount, _ := controller.app.CountRecords(model.DbNameVotes, dbx.HashExp{
			model.VotesFieldArticleId: article.Id,
			model.VotesFieldVoteType:  model.VoteTypeWealth,
		})

		articleRankings = append(articleRankings, map[string]interface{}{
			"rank":            rank + 1,
			"article_id":      article.Id,
			"article_o_id":    article.OId(),
			"title":           article.Title(),
			"preview_content": article.PreviewContent(),
			"view_count":      article.ViewCount(),
			"good_cnt":        article.GoodCnt(),
			"comment_count":   article.CommentCount(),
			"collect_cnt":     article.CollectCnt(),
			"thank_cnt":       article.ThankCnt(),
			"score":           article.Score(),
			"user_id":         user.Id,
			"username":        user.Name(),
			"nickname":        user.Nickname(),
			"avatar":          user.Avatar(),
			"career_votes":    careerCount,
			"romance_votes":   romanceCount,
			"wealth_votes":    wealthCount,
		})
	}

	// 返回所有数据
	return event.JSON(http.StatusOK, map[string]interface{}{
		"champions":    champions,
		"participants": participants,
		"votes": map[string]interface{}{
			"career":  careerVotes,
			"romance": romanceVotes,
			"wealth":  wealthVotes,
		},
		"article_rankings": articleRankings,
	})
}
