package models

import (
	"blog-service/global"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Model
	Title           string `json:"title" gorm:"column:title;varchar(100)"`                     // 文章标题
	Desc            string `json:"desc" gorm:"column:desc;varchar(255)"`                       // 文章简介
	Cover_image_url string `json:"cover_image_url" gorm:"column:cover_image_url;varchar(255)"` // 文章封面图片地址
	CreatedBy       string `json:"created_by" gorm:"column:created_by; varchar(100)"`          // 创建人
	ModifiedBy      string `json:"modified_by" gorm:"varchar(100);column:modified_by"`         // 修改人
	Content         string `json:"content" gorm:"column:content"`
	State           int    `json:"state" gorm:"column:state;tinyint(3)" ` // 状态 0-禁用 1-启用
	// 创建文章标签关联表，这个表主要用于记录文章和标签之间的 1:N 的关联关系
	TagID int `json:"tag_id" gorm:"column:tag_id;int(11)"` // 文章标签ID
	Tag   Tag `json:"tag" gorm:"foreignkey:TagID"`         // 文章标签
}

// 表名
func (a Article) TableName() string {
	return "blog_article"
}

// 新增文章
func AddArticle(data map[string]interface{}) bool {
	err := global.DBEngine.Create(&Article{
		TagID:     data["tag_id"].(int), // 断言
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}).Error
	if err != nil {
		return false
	}
	return true
}

// 由ID判断文章是否存在
func ExistArticleByID(id int) bool {
	var article Article
	err := global.DBEngine.Select("id").Where("id =?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if article.ID > 0 {
		return true
	}
	return false
}

// 获取文章总数
func GetArticleTotal(maps interface{}) (count int64, err error) {
	err = global.DBEngine.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := global.DBEngine.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

// 获取单个文章详情
func GetArticle(id int) (article *Article, err error) {
	err = global.DBEngine.Where("id =?", id).First(&article).Error
	if err != nil {
		return article, err
	}
	err = global.DBEngine.Preload("Tag").Where("id =?", id).First(&article).Error
	if err != nil {
		return article, err
	}
	return article, nil
}

// 更新文章
func UpdateArticle(id int, data interface{}) bool {
	err := global.DBEngine.Model(&Article{}).Where("id =?", id).Updates(data).Error
	if err != nil {
		return false
	}
	return true
}

// 删除文章
func DeleteArticle(id int) bool {
	err := global.DBEngine.Where("id =?", id).Delete(&Article{}).Error
	if err != nil {
		return false
	}
	return true
}

func (article *Article) BeforeCreateA(tx *gorm.DB) error {
	now := time.Now()
	article.CreatedOn = &now
	return nil
}

func (article *Article) BeforeUpdateA(tx *gorm.DB) error {
	now := time.Now()
	article.ModifiedOn = &now
	return nil
}
