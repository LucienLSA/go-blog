package mysql

import (
	"database/sql"

	"github.com/LucienLSA/go-blog/models"
	"go.uber.org/zap"
)

// 数据库创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, 
	author_id, community_id)
	values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content,
		p.AuthorID, p.CommunityID)
	return err
}

// 查询所有帖子列表
func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, author_id, community_id, title
	, content, status, create_time, update_time 
	from post limit ?,?` // 限制只查询出指定条数，分页
	posts = make([]*models.Post, 0, 2)
	if err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no post in database")
			err = nil
		}
	}
	return
}

// 查询分类帖子详情
func GetPostDetailList(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	// fmt.Println(pid)
	sqlStr := `select post_id, author_id, community_id, title
	, content, status, create_time, update_time 
	from post where post_id = ?`
	if err = db.Get(post, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("This post_id is not existing in database")
			err = ErrorInvalidID
		}
	}
	return post, err
}
