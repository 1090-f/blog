package dao

import (
	"blog/internal/model"
	"time"

	"gorm.io/gorm"
)

type PublishedSiteStats struct {
	ArticleCount     int64
	TotalWords       int64
	FirstPublishedAt *time.Time
	LastActivityAt   *time.Time
}

type ArticleListFilter struct {
	CategoryID    uint
	TagID         uint
	Keyword       string
	Status        string
	PublishedOnly bool
	Offset        int
	Limit         int
}

type ArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) *ArticleDAO {
	return &ArticleDAO{db: db}
}

func (d *ArticleDAO) Create(article *model.Article) error {
	return d.db.Create(article).Error
}

func (d *ArticleDAO) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := d.baseQuery().First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (d *ArticleDAO) FindPublishedByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := d.baseQuery().Where("status = ?", "published").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (d *ArticleDAO) ListLatestPublished(limit int) ([]model.Article, error) {
	return d.listPublishedOrdered("created_at DESC", limit)
}

func (d *ArticleDAO) ListPopularPublished(limit int) ([]model.Article, error) {
	return d.listPublishedOrdered("view_count DESC, created_at DESC", limit)
}

func (d *ArticleDAO) List(filter ArticleListFilter) ([]model.Article, int64, error) {
	query := d.filterQuery(filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []model.Article
	if err := d.baseQuery().
		Scopes(func(tx *gorm.DB) *gorm.DB {
			return d.applyFilters(tx, filter)
		}).
		Order("created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (d *ArticleDAO) Update(article *model.Article) error {
	return d.db.Save(article).Error
}

func (d *ArticleDAO) Delete(id uint) error {
	return d.db.Delete(&model.Article{}, id).Error
}

func (d *ArticleDAO) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *ArticleDAO) IncrementViewCount(id uint) error {
	return d.db.Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (d *ArticleDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *ArticleDAO) CountByStatus(status string) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *ArticleDAO) SumViewCount() (int64, error) {
	type result struct {
		Total int64
	}
	var value result
	if err := d.db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0) AS total").Scan(&value).Error; err != nil {
		return 0, err
	}
	return value.Total, nil
}

func (d *ArticleDAO) PublishedSiteStats() (*PublishedSiteStats, error) {
	var stats PublishedSiteStats
	if err := d.db.Model(&model.Article{}).
		Where("status = ?", "published").
		Select("COUNT(*) AS article_count, COALESCE(SUM(CHAR_LENGTH(content)), 0) AS total_words, MIN(created_at) AS first_published_at, MAX(updated_at) AS last_activity_at").
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (d *ArticleDAO) baseQuery() *gorm.DB {
	return d.db.Model(&model.Article{}).Preload("Category").Preload("User").Preload("Tags")
}

func (d *ArticleDAO) filterQuery(filter ArticleListFilter) *gorm.DB {
	return d.applyFilters(d.db.Model(&model.Article{}), filter)
}

func (d *ArticleDAO) applyFilters(tx *gorm.DB, filter ArticleListFilter) *gorm.DB {
	if filter.PublishedOnly {
		tx = tx.Where("status = ?", "published")
	} else if filter.Status != "" {
		tx = tx.Where("status = ?", filter.Status)
	}
	if filter.CategoryID != 0 {
		tx = tx.Where("category_id = ?", filter.CategoryID)
	}
	if filter.TagID != 0 {
		tx = tx.Where("articles.id IN (?)", d.db.Table("article_tags").Select("article_id").Where("tag_id = ?", filter.TagID))
	}
	if filter.Keyword != "" {
		tx = tx.Where("title LIKE ?", "%"+filter.Keyword+"%")
	}
	return tx
}

func (d *ArticleDAO) listPublishedOrdered(order string, limit int) ([]model.Article, error) {
	var articles []model.Article
	if err := d.baseQuery().
		Where("status = ?", "published").
		Order(order).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
