package service

import (
	"bless-activity/model"
	"bless-activity/service/fishpi"
	"fmt"
	"log/slog"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ArticleService struct {
	userMap    *maputil.ConcurrentMap[string, *model.User]
	articleMap *maputil.ConcurrentMap[string, *model.Article]

	app           core.App
	fishpiService *fishpi.Service
}

func NewArticleService(app core.App, fishpiService *fishpi.Service) *ArticleService {

	service := ArticleService{
		userMap:       maputil.NewConcurrentMap[string, *model.User](100),
		articleMap:    maputil.NewConcurrentMap[string, *model.Article](100),
		app:           app,
		fishpiService: fishpiService,
	}
	return &service
}

func (service *ArticleService) Start() {
	service.app.Cron().MustAdd("fetch-article", "*/5 * * * *", service.FetchArticles)
}

func (service *ArticleService) FetchArticles() {

	service.cacheAuthors()
	service.cacheArticles()

	const size = 50
	var page = 1
	for {
		response, err := service.fishpiService.GetApiArticlesTag("福签传情", page, size)
		if err != nil {
			service.app.Logger().Error("爬取文章失败", slog.Any("err", err))
			return
		}
		service.app.Logger().Debug("爬取文章结果", slog.Int("length", len(response.Data.Articles)))

		if len(response.Data.Articles) == 0 {
			break
		}

		// 处理文章 和 作者信息
		for _, article := range response.Data.Articles {
			service.HandleArticle(article)
		}

		if page >= response.Data.Pagination.PaginationPageCount {
			break
		}

		page++
	}

}

func (service *ArticleService) cacheAuthors() {
	var users []*model.User
	if err := service.app.RecordQuery(model.DbNameUsers).All(&users); err != nil {
		return
	}
	for _, user := range users {
		service.userMap.Set(user.OId(), user)
	}
}

func (service *ArticleService) cacheArticles() {
	var articles []*model.Article
	if err := service.app.RecordQuery(model.DbNameArticles).All(&articles); err != nil {
		return
	}
	for _, article := range articles {
		service.articleMap.Set(article.OId(), article)
	}
}

func (service *ArticleService) HandleArticle(responseArticle *fishpi.GetApiArticlesTagResponseArticle) {
	if err := service.HandleAuthor(responseArticle.ArticleAuthor); err != nil {
		return
	}

	article, exist := service.articleMap.Get(responseArticle.OId)
	if exist {
		// 更新文章
		article.SetTitle(responseArticle.ArticleTitle)
		article.SetPreviewContent(responseArticle.ArticlePreviewContent)
		article.SetViewCount(responseArticle.ArticleViewCount)
		article.SetGoodCnt(responseArticle.ArticleGoodCnt)
		article.SetCommentCount(responseArticle.ArticleCommentCount)
		article.SetCollectCnt(responseArticle.ArticleCollectCnt)
		article.SetThankCnt(responseArticle.ArticleThankCnt)
		updatedTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleUpdateTimeStr, time.Local)
		updated, _ := types.ParseDateTime(updatedTime)
		article.SetUpdatedAt(updated)
		if err := service.app.Save(article); err != nil {
			return
		}
		return
	}

	// 创建文章
	user, userExist := service.userMap.Get(responseArticle.ArticleAuthor.OId)
	if !userExist {
		return
	}

	articleCollection, err := service.app.FindCollectionByNameOrId(model.DbNameArticles)
	if err != nil {
		return
	}
	article = model.NewArticleFromCollection(articleCollection)
	article.SetUserId(user.Id)
	article.SetOId(responseArticle.OId)
	article.SetTitle(responseArticle.ArticleTitle)
	article.SetPreviewContent(responseArticle.ArticlePreviewContent)
	article.SetViewCount(responseArticle.ArticleViewCount)
	article.SetGoodCnt(responseArticle.ArticleGoodCnt)
	article.SetCommentCount(responseArticle.ArticleCommentCount)
	article.SetCollectCnt(responseArticle.ArticleCollectCnt)
	article.SetThankCnt(responseArticle.ArticleThankCnt)
	createdTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleCreateTimeStr, time.Local)
	created, _ := types.ParseDateTime(createdTime)
	article.SetCreatedAt(created)
	updatedTime, _ := time.ParseInLocation(time.DateTime, responseArticle.ArticleUpdateTimeStr, time.Local)
	updated, _ := types.ParseDateTime(updatedTime)
	article.SetUpdatedAt(updated)
	if err = service.app.Save(article); err != nil {
		return
	}
	service.articleMap.Set(article.OId(), article)
}

func (service *ArticleService) HandleAuthor(author *fishpi.GetApiArticlesTagResponseArticleAuthor) error {
	user, exist := service.userMap.Get(author.OId)
	if exist {
		// 更新用户
		user.SetName(author.UserName)
		user.SetNickname(author.UserNickname)
		user.SetAvatar(author.UserAvatarURL)
		if err := service.app.Save(user); err != nil {
			return err
		}
		return nil
	}
	// 创建用户
	userCollection, err := service.app.FindCollectionByNameOrId(model.DbNameUsers)
	if err != nil {
		return err
	}
	user = model.NewUserFromCollection(userCollection)
	user.SetEmail(fmt.Sprintf("%s@fishpi.cn", author.OId))
	user.SetEmailVisibility(true)
	user.SetVerified(true)
	user.SetOId(author.OId)
	user.SetName(author.UserName)
	user.SetNickname(author.UserNickname)
	user.SetAvatar(author.UserAvatarURL)
	user.SetRandomPassword()
	if err = service.app.Save(user); err != nil {
		return err
	}
	service.userMap.Set(user.OId(), user)
	return nil
}
