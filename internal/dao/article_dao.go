package dao

import (
	"blog/internal/model"
	"time"

	"gorm.io/gorm"
)

// PublishedSiteStats 已发布文章的站点聚合统计。
type PublishedSiteStats struct {
	ArticleCount     int64
	TotalWords       int64
	FirstPublishedAt *time.Time
	LastActivityAt   *time.Time
}

// ArticleListFilter 文章分页查询的筛选条件。
type ArticleListFilter struct {
	CategoryID    uint
	TagID         uint
	Keyword       string
	Status        string
	PublishedOnly bool
	Offset        int
	Limit         int
}

// ArticleDAO 文章表的数据访问层，支持复杂筛选和聚合查询。
type ArticleDAO struct {
	db *gorm.DB
}

// NewArticleDAO 创建并初始化文章数据访问层实例。
func NewArticleDAO(db *gorm.DB) *ArticleDAO {
	return &ArticleDAO{db: db}
}

// Create 新增文章记录到数据库。
func (d *ArticleDAO) Create(article *model.Article) error {
	return d.db.Create(article).Error
}

// FindByID 按主键查询文章，预加载分类、用户和标签关联。
func (d *ArticleDAO) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := d.baseQuery().First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// FindPublishedByID 按主键查询已发布文章，预加载分类、用户和标签关联。
func (d *ArticleDAO) FindPublishedByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := d.baseQuery().Where("status = ?", "published").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// ListLatestPublished 按创建时间倒序查询最新的已发布文章列表。
func (d *ArticleDAO) ListLatestPublished(limit int) ([]model.Article, error) {
	return d.listPublishedOrdered("created_at DESC", limit)
}

// ListPopularPublished 按浏览量和创建时间倒序查询热门已发布文章列表。
func (d *ArticleDAO) ListPopularPublished(limit int) ([]model.Article, error) {
	return d.listPublishedOrdered("view_count DESC, created_at DESC", limit)
}

// List 按筛选条件分页查询文章列表，先统计总数再查询当前页数据。
func (d *ArticleDAO) List(filter ArticleListFilter) ([]model.Article, int64, error) {
	query := d.filterQuery(filter)

	// 先统计符合筛选条件的总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询当前页的文章数据，按创建时间倒序
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

// Update 更新文章记录（全量保存）。
func (d *ArticleDAO) Update(article *model.Article) error {
	return d.db.Save(article).Error
}

// Delete 按主键删除文章记录。
func (d *ArticleDAO) Delete(id uint) error {
	return d.db.Delete(&model.Article{}, id).Error
}

// CountByCategoryID 统计指定分类下的文章数量。
func (d *ArticleDAO) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// IncrementViewCount 原子递增指定文章的浏览量，使用 SQL 表达式避免并发问题。
func (d *ArticleDAO) IncrementViewCount(id uint) error {
	return d.db.Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// CountAll 统计所有文章的总数量。
func (d *ArticleDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus 统计指定状态的文章数量。
func (d *ArticleDAO) CountByStatus(status string) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Article{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SumViewCount 汇总所有文章的浏览量总和，无数据时返回 0。
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

// PublishedSiteStats 聚合查询已发布文章的站点统计数据（文章数、总字数、首次发布时间、最近更新时间）。
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

// baseQuery 构建基础查询，预加载分类、用户和标签关联。
func (d *ArticleDAO) baseQuery() *gorm.DB {
	return d.db.Model(&model.Article{}).Preload("Category").Preload("User").Preload("Tags")
}

// filterQuery 构建带筛选条件的查询（不含关联预加载）。
func (d *ArticleDAO) filterQuery(filter ArticleListFilter) *gorm.DB {
	return d.applyFilters(d.db.Model(&model.Article{}), filter)
}

// applyFilters 将筛选条件链式应用到查询中，支持状态、分类、标签和关键词筛选。
func (d *ArticleDAO) applyFilters(tx *gorm.DB, filter ArticleListFilter) *gorm.DB {
	// 状态筛选：PublishedOnly 为 true 时只查已发布，否则按指定状态筛选
	if filter.PublishedOnly {
		tx = tx.Where("status = ?", "published")
	} else if filter.Status != "" {
		tx = tx.Where("status = ?", filter.Status)
	}
	// 分类筛选
	if filter.CategoryID != 0 {
		tx = tx.Where("category_id = ?", filter.CategoryID)
	}
	// 标签筛选：通过 article_tags 关联表进行子查询
	if filter.TagID != 0 {
		tx = tx.Where("articles.id IN (?)", d.db.Table("article_tags").Select("article_id").Where("tag_id = ?", filter.TagID))
	}
	// 关键词模糊搜索（仅匹配标题）
	if filter.Keyword != "" {
		tx = tx.Where("title LIKE ?", "%"+filter.Keyword+"%")
	}
	return tx
}

// listPublishedOrdered 内部方法：按指定排序方式查询已发布文章列表。
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
