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

// 查询所有帖子
func GetPost() (postList []*models.PostList, err error) {
	sqlStr := `select post_id, title, status from post `
	if err = db.Select(&postList, sqlStr); err != nil {
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
			err = ErrorInvalidID
		}
	}
	return post, err
}
