package application

import (
	"bless-activity/service"
	"bless-activity/service/fishpi"
	"log/slog"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

/*
	活动：互送祝福、互动排行、激情博饼
	页面：排行榜（可排序-互动、财运、姻缘、事业）博饼奖励 博饼参与 博饼记录
	后端：定时统计文章信息 为作者创建账号 更新文章排行
    作者登陆账号 进行赠送祝福和博饼


*/

type Application struct {
	app *pocketbase.PocketBase

	articleService *service.ArticleService
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

	fishPiService, err := fishpi.NewService(event.App)
	if err != nil {
		event.App.Logger().Error("创建fishPi Service失败", slog.Any("err", err))
		return err
	}

	// 文章爬取服务
	application.articleService = service.NewArticleService(event.App, fishPiService)
	go application.articleService.FetchArticles()

	// 注册路由
	application.registerRoutes()

	event.App.Logger().Debug("初始化完成")
	return nil
}

func (application *Application) registerRoutes() {
	application.app.OnServe().BindFunc(func(event *core.ServeEvent) error {

		return event.Next()
	})
}
