package service

import (
	"bless-activity/service/fishpi"

	"github.com/pocketbase/pocketbase/core"
)

type ArticleService struct {
	app core.App

	fishpiService *fishpi.Service
}

func NewArticleService(app core.App, fishpiService *fishpi.Service) *ArticleService {
	service := ArticleService{
		app:           app,
		fishpiService: fishpiService,
	}
	return &service
}

func (service *ArticleService) FetchArticles() {
	const size = 50
	var page = 1
	for {
		response, err := service.fishpiService.GetApiArticlesTag("福签传情", page, size)
		if err != nil {
			return
		}

		if len(response.Data.Articles) == 0 {
			break
		}
		if page >= response.Data.Pagination.PaginationPageCount {
			break
		}

		// 处理文章 和 作者信息
		for _, article := range response.Data.Articles {

		}

		page++
	}

}
