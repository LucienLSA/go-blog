package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
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
func (dao *PostDao) CreatePost(p *models.Post) (err error) {
	// sqlStr := `insert into post(post_id, title, content,
	// author_id, community_id)
	// values(?,?,?,?,?)`

	// _, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content,
	// 	p.AuthorID, p.CommunityID)
	// return err
	return
}

// // 查询所有帖子列表
// func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
// 	sqlStr := `select post_id, author_id, community_id, title
// 	, content, status, create_time, update_time
// 	from post
// 	ORDER BY update_time
// 	DESC
// 	limit ?,?`
// 	// 限制只查询出指定条数，分页
// 	posts = make([]*models.Post, 0, 2)
// 	if err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize); err != nil {
// 		if err == sql.ErrNoRows {
// 			zap.L().Warn("there is no post in database")
// 			err = nil
// 		}
// 	}
// 	return
// }

// // 根据帖子id查询单个分类帖子详情
// func GetPostDetailList(pid int64) (post *models.Post, err error) {
// 	post = new(models.Post)
// 	// fmt.Println(pid)
// 	sqlStr := `select post_id, author_id, community_id, title
// 	, content, status, create_time, update_time
// 	from post where post_id = ?`
// 	if err = db.Get(post, sqlStr, pid); err != nil {
// 		if err == sql.ErrNoRows {
// 			zap.L().Warn("This post_id is not existing in database")
// 			err = ErrorInvalidID
// 		}
// 	}
// 	return post, err
// }

// // 根据给定的ids列表查询帖子
// func SearchPostListByIDs(ids []string) (postList []*models.Post, err error) {
// 	sqlStr := `select post_id, author_id, community_id, title
// 	, content, status, create_time, update_time
// 	from post where post_id in (?)
// 	order by FIND_IN_SET(post_id,?)`
// 	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
// 	if err != nil {
// 		return nil, err
// 	}
// 	query = db.Rebind(query)
// 	err = db.Select(&postList, query, args...)
// 	return
// }
