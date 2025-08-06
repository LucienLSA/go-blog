package mysql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostDao struct {
	*gorm.DB
}

func NewPostDao(ctx context.Context) *PostDao {
	return &PostDao{NewDBClient(ctx)}
}

func NewPostDaoByDB(db *gorm.DB) *PostDao {
	return &PostDao{db}
}

// 数据库创建帖子
func (dao *PostDao) CreatePost(post *models.Post) (err error) {
	// sqlStr := `insert into post(post_id, title, content,
	// author_id, community_id)
	// values(?,?,?,?,?)`

	// _, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content,
	// 	p.AuthorID, p.CommunityID)
	// return err
	return dao.DB.Model(&models.Post{}).Create(&post).Error
}

// 查询所有帖子列表
func (dao *PostDao) GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	// sqlStr := `select post_id, author_id, community_id, title
	// , content, status, create_time, update_time
	// from post
	// ORDER BY update_time
	// DESC
	// limit ?,?`
	// // 限制只查询出指定条数，分页
	// posts = make([]*models.Post, 0, 2)
	// if err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		zap.L().Warn("there is no post in database")
	// 		err = nil
	// 	}
	// }
	offset := (pageNum - 1) * pageSize
	err = dao.DB.Model(&models.Post{}).
		Order("updated_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&posts).Error
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no post in database")
		return nil, err
	}
	return
}

// 根据帖子id查询单个分类帖子详情
func (dao *PostDao) GetPostDetailList(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	err = dao.DB.Model(&models.Post{}).Where("post_id=?", pid).First(post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return post, nil
}

// 根据给定的ids列表查询帖子
func (dao *PostDao) SearchPostListByIDs(ids []string) (postList []*models.Post, err error) {
	// sqlStr := `select post_id, author_id, community_id, title
	// , content, status, create_time, update_time
	// from post where post_id in (?)
	// order by FIND_IN_SET(post_id,?)`
	// query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	// if err != nil {
	// 	return nil, err
	// }
	// query = db.Rebind(query)
	// err = db.Select(&postList, query, args...)
	// 将 ids 转换为字符串，用于 FIND_IN_SET
	idsStr := strings.Join(ids, ",")
	err = dao.DB.Model(&models.Post{}).Where("post_id IN (?)", ids).
		Order(gorm.Expr("FIND_IN_SET(post_id, ?)", idsStr)).
		Find(&postList).Error
	return
}

// CreatePostWithReview 创建帖子并设置为待审核状态
func (dao *PostDao) CreatePostWithReview(post *models.Post) error {
	// 设置初始审核状态
	post.ReviewStatus = "pending"
	post.ReviewScore = 0.0
	
	return dao.DB.Model(&models.Post{}).Create(post).Error
}

// UpdatePostReviewStatus 更新帖子审核状态
func (dao *PostDao) UpdatePostReviewStatus(postID int64, reviewStatus string, reviewResult string) error {
	updates := map[string]interface{}{
		"review_status": reviewStatus,
		"review_result": reviewResult,
	}
	
	if reviewStatus == "approved" || reviewStatus == "rejected" {
		now := time.Now().Format("2006-01-02 15:04:05")
		updates["reviewed_at"] = &now
	}
	
	return dao.DB.Model(&models.Post{}).Where("post_id = ?", postID).Updates(updates).Error
}

// GetPostsByReviewStatus 根据审核状态获取帖子列表
func (dao *PostDao) GetPostsByReviewStatus(status string, pageNum, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (pageNum - 1) * pageSize
	
	err := dao.DB.Model(&models.Post{}).
		Where("review_status = ?", status).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&posts).Error
	
	if err == sql.ErrNoRows {
		zap.L().Warn("no posts found with review status", zap.String("status", status))
		return nil, nil
	}
	
	return posts, err
}

// GetPostWithReviewStatus 获取帖子详情（包含审核状态）
func (dao *PostDao) GetPostWithReviewStatus(postID int64) (*models.Post, error) {
	var post models.Post
	err := dao.DB.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}
