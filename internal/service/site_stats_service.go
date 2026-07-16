package service

import (
	"blog/internal/dao"
	"time"
)

// PublicSiteStats 站点公开统计数据（文章数、分类数、标签数、总字数等）。
type PublicSiteStats struct {
	ArticleCount     int64      `json:"articleCount"`
	CategoryCount    int64      `json:"categoryCount"`
	TagCount         int64      `json:"tagCount"`
	TotalWords       int64      `json:"totalWords"`
	FirstPublishedAt *time.Time `json:"firstPublishedAt"`
	LastActivityAt   *time.Time `json:"lastActivityAt"`
}

// SiteStatsArticleStore 站点统计所需的文章聚合操作抽象。
type SiteStatsArticleStore interface {
	PublishedSiteStats() (*dao.PublishedSiteStats, error)
}

// SiteStatsCategoryStore 站点统计所需的分类计数操作抽象。
type SiteStatsCategoryStore interface {
	CountAll() (int64, error)
}

// SiteStatsTagStore 站点统计所需的标签计数操作抽象。
type SiteStatsTagStore interface {
	CountAll() (int64, error)
}

// SiteStatsService 站点公开统计业务逻辑层。
type SiteStatsService struct {
	articleStore  SiteStatsArticleStore
	categoryStore SiteStatsCategoryStore
	tagStore      SiteStatsTagStore
}

// NewSiteStatsService 创建并初始化站点统计实例。
func NewSiteStatsService(articleStore SiteStatsArticleStore, categoryStore SiteStatsCategoryStore, tagStore SiteStatsTagStore) *SiteStatsService {
	return &SiteStatsService{articleStore: articleStore, categoryStore: categoryStore, tagStore: tagStore}
}

// Get 获取站点公开统计数据，包括已发布文章数、分类数、标签数、总字数、首次发布时间和最近活动时间。
func (s *SiteStatsService) Get() (*PublicSiteStats, error) {
	// 查询已发布文章的聚合数据（文章数、总字数、首次/最近发布时间）
	articleStats, err := s.articleStore.PublishedSiteStats()
	if err != nil {
		return nil, err
	}
	// 查询所有分类数量
	categoryCount, err := s.categoryStore.CountAll()
	if err != nil {
		return nil, err
	}
	// 查询所有标签数量
	tagCount, err := s.tagStore.CountAll()
	if err != nil {
		return nil, err
	}

	// 组装公开统计数据结构体返回
	return &PublicSiteStats{
		ArticleCount:     articleStats.ArticleCount,
		CategoryCount:    categoryCount,
		TagCount:         tagCount,
		TotalWords:       articleStats.TotalWords,
		FirstPublishedAt: articleStats.FirstPublishedAt,
		LastActivityAt:   articleStats.LastActivityAt,
	}, nil
}
