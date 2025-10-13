package application

import (
	"bless-activity/controller"
	"bless-activity/service"
	"bless-activity/service/fishpi"
	"log/slog"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

/*
	活动：互送祝福、互动排行、激情博饼
	页面：排行榜（可排序-互动、财运、姻缘、事业）博饼奖励 博饼参与 博饼记录
	后端：定时统计文章信息 为作者创建账号 更新文章排行
    作者登陆账号 进行赠送祝福和博饼

    /login 登陆
    /callback 回调

    /me name nickname avatar base_times extra_times rest_times article_id
    /rank 排行榜
    /awards 博饼奖励内容 剩余数量
    /play 博饼参与
    /histories 博饼记录
    /vote 投票
    /vote/list 投票记录



*/

type Application struct {
	app *pocketbase.PocketBase

	fishPiService  *fishpi.Service
	articleService *service.ArticleService

	fishPiController   *controller.FishPiController
	userController     *controller.UserController
	mooncakeController *controller.MooncakeController
	voteController     *controller.VoteController
}

func NewApp() *Application {
	app := pocketbase.New()

	application := &Application{
		app: app,
	}

	return application
}

func (application *Application) Start() error {

	// 初始化
	application.app.OnBootstrap().BindFunc(func(event *core.BootstrapEvent) error {

		if err := event.Next(); err != nil {
			return err
		}

		return application.init(event)
	})

	return application.app.Start()
}

func (application *Application) init(event *core.BootstrapEvent) error {
	event.App.Logger().Debug("初始化程序")

	var err error
	if application.fishPiService, err = fishpi.NewService(event.App); err != nil {
		event.App.Logger().Error("创建fishPi Service失败", slog.Any("err", err))
		return err
	}

	// 文章爬取服务
	application.articleService = service.NewArticleService(event.App, application.fishPiService)
	go application.articleService.FetchArticles()

	// 注册路由
	application.app.OnServe().BindFunc(application.registerRoutes)

	event.App.Logger().Debug("初始化完成")
	return nil
}

func (application *Application) registerRoutes(event *core.ServeEvent) error {

	event.Router.Bind(
		&hook.Handler[*core.RequestEvent]{
			Id:       "cookie_to_header",
			Priority: apis.DefaultRateLimitMiddlewarePriority - 999,
			Func: func(event *core.RequestEvent) error {
				if event.Request.Header.Get("Authorization") != "" {
					return event.Next()
				}

				cookie, err := event.Request.Cookie("token")
				if err != nil {
					return event.Next()
				}
				if err = cookie.Valid(); err != nil {
					return event.Next()
				}
				event.Request.Header.Set("Authorization", cookie.Value)

				return event.Next()
			},
		},
		&hook.Handler[*core.RequestEvent]{
			Id:       "query_to_header",
			Priority: apis.DefaultRateLimitMiddlewarePriority - 998,
			Func: func(event *core.RequestEvent) error {
				if event.Request.Header.Get("Authorization") != "" {
					return event.Next()
				}

				queryAuth := event.Request.URL.Query().Get("authorization")
				if queryAuth == "" {
					return event.Next()
				}
				event.Request.Header.Set("Authorization", queryAuth)

				return event.Next()
			},
		},
	)

	application.fishPiController = controller.NewFishPiController(event)
	application.userController = controller.NewUserController(event)
	application.mooncakeController = controller.NewMooncakeController(event)
	application.voteController = controller.NewVoteController(event)

	event.Router.GET("/test", func(e *core.RequestEvent) error {
		return e.String(http.StatusOK, "test")
	})
	event.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

	event.Router.GET("/test1", func(e *core.RequestEvent) error {
		return e.String(http.StatusOK, "test1")
	})

	return event.Next()
}
