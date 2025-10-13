package controller

import (
	"bless-activity/model"
	"log/slog"
	"net/http"
	"time"

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
	// 后端登出，清除 token cookie 并重定向到首页
	group.GET("/logout", controller.Logout)
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

	totalTimes := model.DefaultMooncakeGamblingTimes + article.ThankCnt()
	// 限制最大次数为20次
	if totalTimes > model.MaxMooncakeGamblingTimes {
		totalTimes = model.MaxMooncakeGamblingTimes
	}
	restTimes := totalTimes - int(drawTimes)

	return event.JSON(http.StatusOK, map[string]any{
		"id":                              user.Id,
		"o_id":                            user.OId(),
		"name":                            user.Name(),
		"nickname":                        user.Nickname(),
		"avatar":                          user.Avatar(),
		"article_id":                      article.Id,
		"article_o_id":                    article.OId(),
		"article_title":                   article.Title(),
		"article_thank_cnt":               article.ThankCnt(),
		"default_mooncake_gambling_times": model.DefaultMooncakeGamblingTimes,
		"max_mooncake_gambling_times":     model.MaxMooncakeGamblingTimes,
		"draw_times":                      drawTimes,
		"rest_times":                      restTimes,
	})
}

// Logout 清除 token cookie 并重定向到首页
func (controller *UserController) Logout(event *core.RequestEvent) error {
	// 设置过期 cookie 清除客户端 token（兼容HttpOnly和非HttpOnly情形）
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  (func() time.Time { t := time.Unix(0, 0); return t })(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	event.SetCookie(cookie)

	// 也尝试清除带点域名的 cookie（部分场景需要）
	// 注意: Go 的 http.SetCookie 无法直接带 domain 设置为 .example.com 进行删除
	// 如果需要，请在前端或部署环境中统一处理 domain

	return event.Redirect(http.StatusFound, "/")
}
