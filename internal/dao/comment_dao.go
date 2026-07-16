package dao

import (
	"blog/internal/model"
	"strings"

	"gorm.io/gorm"
)

// CommentDAO 评论表的数据访问层。
type CommentDAO struct {
	db *gorm.DB
}

// CommentListFilter 管理端评论分页查询的筛选条件。
type CommentListFilter struct {
	Keyword   string
	ArticleID uint
	Status    *int8
	Offset    int
	Limit     int
}

// NewCommentDAO 创建并初始化评论数据访问层实例。
func NewCommentDAO(db *gorm.DB) *CommentDAO {
	return &CommentDAO{db: db}
}

// Create 新增评论记录到数据库。
func (d *CommentDAO) Create(comment *model.Comment) error {
	return d.db.Create(comment).Error
}

// FindByID 按主键查询评论，预加载文章、用户和被回复评论及其用户关联。
func (d *CommentDAO) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := d.baseQuery().First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByArticleID 查询指定文章下所有已审核通过的评论，按创建时间正序排列。
func (d *CommentDAO) ListByArticleID(articleID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := d.baseQuery().
		Where("article_id = ? AND status = ?", articleID, 1).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// ListAdmin 按筛选条件分页查询管理端评论列表，先统计总数再查询当前页数据。
func (d *CommentDAO) ListAdmin(filter CommentListFilter) ([]model.Comment, int64, error) {
	query := applyCommentFilters(d.adminQuery(), filter)

	// 先统计符合筛选条件的总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询当前页的评论数据，按创建时间倒序
	var comments []model.Comment
	if err := applyCommentFilters(d.adminQuery(), filter).
		Order("comments.created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateStatus 更新指定评论的审核状态字段。
func (d *CommentDAO) UpdateStatus(id uint, status int8) error {
	return d.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

// CountAll 统计所有评论的总数量。
func (d *CommentDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Comment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Delete 按主键删除评论记录。
func (d *CommentDAO) Delete(id uint) error {
	return d.db.Delete(&model.Comment{}, id).Error
}

// baseQuery 构建基础查询，预加载文章、用户和被回复评论及其用户关联。
func (d *CommentDAO) baseQuery() *gorm.DB {
	return d.db.Model(&model.Comment{}).
		Preload("Article").
		Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User")
}

// adminQuery 构建管理端评论查询，额外 JOIN 用户表和文章表以支持跨表筛选。
func (d *CommentDAO) adminQuery() *gorm.DB {
	return d.baseQuery().
		Joins("LEFT JOIN users ON users.id = comments.user_id").
		Joins("JOIN articles ON articles.id = comments.article_id")
}

// applyCommentFilters 将筛选条件链式应用到查询中，支持按文章 ID、状态和关键词筛选。
// 关键词匹配评论内容、游客昵称/邮箱、用户名/昵称和文章标题。
func applyCommentFilters(query *gorm.DB, filter CommentListFilter) *gorm.DB {
	// 按文章 ID 筛选
	if filter.ArticleID != 0 {
		query = query.Where("comments.article_id = ?", filter.ArticleID)
	}
	// 按审核状态筛选（指针类型，nil 表示不筛选）
	if filter.Status != nil {
		query = query.Where("comments.status = ?", *filter.Status)
	}
	// 关键词模糊搜索：匹配评论内容、游客信息、用户信息和文章标题
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"(comments.content LIKE ? OR comments.guest_name LIKE ? OR comments.guest_email LIKE ? OR users.username LIKE ? OR users.nickname LIKE ? OR articles.title LIKE ?)",
			like, like, like, like, like, like,
		)
	}
	return query
}
