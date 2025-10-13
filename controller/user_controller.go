package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type UserController struct {
	event *core.ServeEvent
	app   core.App

	logger *slog.Logger
}

func NewUserController(event *core.ServeEvent) *UserController {
	logger := event.App.Logger().With(
		slog.String("controller", "user"),
	)

	controller := &UserController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *UserController) registerRoutes() {
	group := controller.event.Router.Group("/user")
	group.GET("/me", controller.GetMe).BindFunc(
		controller.CheckLogin,
	)
}

func (controller *UserController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *UserController) CheckLogin(event *core.RequestEvent) error {
	if event.Auth == nil {
		return event.UnauthorizedError("未登录", nil)
	}
	if event.HasSuperuserAuth() {
		return event.ForbiddenError("请登录普通用户账号", nil)
	}
	return event.Next()
}

func (controller *UserController) GetMe(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("get_me")

	user := model.NewUser(event.Auth)

	article := new(model.Article)
	if err := controller.app.RecordQuery(model.DbNameArticles).Where(dbx.HashExp{model.ArticlesFieldUserId: user.Id}).OrderBy(model.ArticlesFieldCreatedAt + " desc").One(article); err != nil {
		logger.Error("查找最新文章失败", slog.Any("err", err))
		return event.InternalServerError("查找最新文章失败", err)
	}

	drawTimes, drawTimesErr := controller.app.CountRecords(model.DbNameHistories, dbx.HashExp{model.HistoriesFieldUserId: user.Id})
	if drawTimesErr != nil {
		logger.Error("查找抽奖次数失败", slog.Any("err", drawTimesErr))
		return event.InternalServerError("查找抽奖次数失败", drawTimesErr)
	}

	restTimes := model.DefaultMooncakeGamblingTimes + article.ThankCnt() - int(drawTimes)

	return event.JSON(http.StatusOK, map[string]any{
		"id":                              user.Id,
		"o_id":                            user.OId(),
		"name":                            user.Name(),
		"nickname":                        user.Nickname(),
		"avatar":                          user.Avatar(),
		"article_id":                      article.Id,
		"article_title":                   article.Title(),
		"article_thank_cnt":               article.ThankCnt(),
		"default_mooncake_gambling_times": model.DefaultMooncakeGamblingTimes,
		"draw_times":                      drawTimes,
		"rest_times":                      restTimes,
	})
}
