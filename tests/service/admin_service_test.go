package service_test

import (
	"testing"

	"blog/internal/service"
)

type fakeAdminStatsArticleStore struct{}

func (f *fakeAdminStatsArticleStore) CountAll() (int64, error) { return 12, nil }
func (f *fakeAdminStatsArticleStore) CountByStatus(status string) (int64, error) {
	if status == service.ArticleStatusPublished {
		return 8, nil
	}
	return 4, nil
}
func (f *fakeAdminStatsArticleStore) SumViewCount() (int64, error) { return 345, nil }

type fakeAdminStatsCategoryStore struct{}

func (f *fakeAdminStatsCategoryStore) CountAll() (int64, error) { return 3, nil }

type fakeAdminStatsUserStore struct{}

func (f *fakeAdminStatsUserStore) CountAll() (int64, error) { return 6, nil }

func TestDashboardReturnsAggregatedStats(t *testing.T) {
	svc := service.NewAdminService(
		&fakeAdminStatsArticleStore{},
		&fakeAdminStatsCategoryStore{},
		&fakeAdminStatsUserStore{},
	)

	stats, err := svc.Dashboard()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if stats.ArticleCount != 12 || stats.PublishedCount != 8 || stats.DraftCount != 4 {
		t.Fatalf("unexpected article stats: %+v", stats)
	}
	if stats.CategoryCount != 3 || stats.UserCount != 6 || stats.TotalViews != 345 {
		t.Fatalf("unexpected dashboard stats: %+v", stats)
	}
}
