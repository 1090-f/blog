package service

import (
	"blog/internal/dao"
	"time"
)

type PublicSiteStats struct {
	ArticleCount     int64      `json:"articleCount"`
	CategoryCount    int64      `json:"categoryCount"`
	TagCount         int64      `json:"tagCount"`
	TotalWords       int64      `json:"totalWords"`
	FirstPublishedAt *time.Time `json:"firstPublishedAt"`
	LastActivityAt   *time.Time `json:"lastActivityAt"`
}

type SiteStatsArticleStore interface {
	PublishedSiteStats() (*dao.PublishedSiteStats, error)
}

type SiteStatsCategoryStore interface {
	CountAll() (int64, error)
}

type SiteStatsTagStore interface {
	CountAll() (int64, error)
}

type SiteStatsService struct {
	articleStore  SiteStatsArticleStore
	categoryStore SiteStatsCategoryStore
	tagStore      SiteStatsTagStore
}

func NewSiteStatsService(articleStore SiteStatsArticleStore, categoryStore SiteStatsCategoryStore, tagStore SiteStatsTagStore) *SiteStatsService {
	return &SiteStatsService{articleStore: articleStore, categoryStore: categoryStore, tagStore: tagStore}
}

func (s *SiteStatsService) Get() (*PublicSiteStats, error) {
	articleStats, err := s.articleStore.PublishedSiteStats()
	if err != nil {
		return nil, err
	}
	categoryCount, err := s.categoryStore.CountAll()
	if err != nil {
		return nil, err
	}
	tagCount, err := s.tagStore.CountAll()
	if err != nil {
		return nil, err
	}

	return &PublicSiteStats{
		ArticleCount:     articleStats.ArticleCount,
		CategoryCount:    categoryCount,
		TagCount:         tagCount,
		TotalWords:       articleStats.TotalWords,
		FirstPublishedAt: articleStats.FirstPublishedAt,
		LastActivityAt:   articleStats.LastActivityAt,
	}, nil
}
