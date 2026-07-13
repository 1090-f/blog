package dao

import (
	"blog/internal/model"
	"strings"

	"gorm.io/gorm"
)

type CommentDAO struct {
	db *gorm.DB
}

type CommentListFilter struct {
	Keyword   string
	ArticleID uint
	Status    *int8
	Offset    int
	Limit     int
}

func NewCommentDAO(db *gorm.DB) *CommentDAO {
	return &CommentDAO{db: db}
}

func (d *CommentDAO) Create(comment *model.Comment) error {
	return d.db.Create(comment).Error
}

func (d *CommentDAO) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := d.baseQuery().First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

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

func (d *CommentDAO) ListAdmin(filter CommentListFilter) ([]model.Comment, int64, error) {
	query := applyCommentFilters(d.adminQuery(), filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

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

func (d *CommentDAO) UpdateStatus(id uint, status int8) error {
	return d.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (d *CommentDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Comment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *CommentDAO) Delete(id uint) error {
	return d.db.Delete(&model.Comment{}, id).Error
}

func (d *CommentDAO) baseQuery() *gorm.DB {
	return d.db.Model(&model.Comment{}).
		Preload("Article").
		Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User")
}

func (d *CommentDAO) adminQuery() *gorm.DB {
	return d.baseQuery().
		Joins("LEFT JOIN users ON users.id = comments.user_id").
		Joins("JOIN articles ON articles.id = comments.article_id")
}

func applyCommentFilters(query *gorm.DB, filter CommentListFilter) *gorm.DB {
	if filter.ArticleID != 0 {
		query = query.Where("comments.article_id = ?", filter.ArticleID)
	}
	if filter.Status != nil {
		query = query.Where("comments.status = ?", *filter.Status)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"(comments.content LIKE ? OR comments.guest_name LIKE ? OR comments.guest_email LIKE ? OR users.username LIKE ? OR users.nickname LIKE ? OR articles.title LIKE ?)",
			like, like, like, like, like, like,
		)
	}
	return query
}
