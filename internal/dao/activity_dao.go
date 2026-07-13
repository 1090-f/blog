package dao

import (
	"blog/internal/model"
	"time"

	"gorm.io/gorm"
)

type DailyActivityCount struct {
	Day   time.Time `gorm:"column:activity_day"`
	Count int64     `gorm:"column:count"`
}

type ActivityDAO struct {
	db *gorm.DB
}

func NewActivityDAO(db *gorm.DB) *ActivityDAO {
	return &ActivityDAO{db: db}
}

func (d *ActivityDAO) PublishedArticlesByDay(start, end time.Time) ([]DailyActivityCount, error) {
	var rows []DailyActivityCount
	err := d.db.Model(&model.Article{}).
		Select("DATE(created_at) AS activity_day, COUNT(*) AS count").
		Where("status = ? AND created_at >= ? AND created_at < ?", "published", start, end).
		Group("DATE(created_at)").
		Scan(&rows).Error
	return rows, err
}

func (d *ActivityDAO) ApprovedCommentsByDay(start, end time.Time) ([]DailyActivityCount, error) {
	var rows []DailyActivityCount
	err := d.db.Model(&model.Comment{}).
		Select("DATE(comments.created_at) AS activity_day, COUNT(*) AS count").
		Joins("JOIN articles ON articles.id = comments.article_id").
		Where("comments.status = ? AND articles.status = ? AND comments.created_at >= ? AND comments.created_at < ?", 1, "published", start, end).
		Group("DATE(comments.created_at)").
		Scan(&rows).Error
	return rows, err
}
