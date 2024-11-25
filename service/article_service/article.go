package article_service

import (
	"blog-service/models"
	"blog-service/pkg/gredis"
	"blog-service/pkg/logging"
	"blog-service/service/cache_service"
	"encoding/json"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.LogrusObj.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) Count() (int64, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}

// 获取所有的文章
func (a *Article) GetAll() ([]*models.Article, error) {
	var articles []*models.Article
	var cacheArticles []*models.Article

	cache := cache_service.Article{
		TagID:    a.TagID,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
		State:    a.State,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.LogrusObj.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, articles, 3600)
	return articles, nil
}
