package service_test

import (
	"testing"
	"time"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	"gorm.io/gorm"
)

type fakeArticleStore struct {
	listFn        func(filter dao.ArticleListFilter) ([]model.Article, int64, error)
	findByIDFn    func(id uint) (*model.Article, error)
	findPubByIDFn func(id uint) (*model.Article, error)
	createFn      func(article *model.Article) error
	updateFn      func(article *model.Article) error
	deleteFn      func(id uint) error
	countByCatFn  func(categoryID uint) (int64, error)
	incrementFn   func(id uint) error
	latestFn      func(limit int) ([]model.Article, error)
	popularFn     func(limit int) ([]model.Article, error)
}

func (f *fakeArticleStore) Create(article *model.Article) error      { return f.createFn(article) }
func (f *fakeArticleStore) FindByID(id uint) (*model.Article, error) { return f.findByIDFn(id) }
func (f *fakeArticleStore) FindPublishedByID(id uint) (*model.Article, error) {
	return f.findPubByIDFn(id)
}
func (f *fakeArticleStore) ListLatestPublished(limit int) ([]model.Article, error) {
	return f.latestFn(limit)
}
func (f *fakeArticleStore) ListPopularPublished(limit int) ([]model.Article, error) {
	return f.popularFn(limit)
}
func (f *fakeArticleStore) List(filter dao.ArticleListFilter) ([]model.Article, int64, error) {
	return f.listFn(filter)
}
func (f *fakeArticleStore) Update(article *model.Article) error { return f.updateFn(article) }
func (f *fakeArticleStore) Delete(id uint) error                { return f.deleteFn(id) }
func (f *fakeArticleStore) CountByCategoryID(categoryID uint) (int64, error) {
	return f.countByCatFn(categoryID)
}
func (f *fakeArticleStore) IncrementViewCount(id uint) error { return f.incrementFn(id) }

type fakeCategoryStore struct {
	findByIDFn func(id uint) (*model.Category, error)
}

func (f *fakeCategoryStore) FindByID(id uint) (*model.Category, error) {
	return f.findByIDFn(id)
}

type fakeCommentReadStore struct {
	listByArticleFn func(articleID uint) ([]model.Comment, error)
}

func (f *fakeCommentReadStore) ListByArticleID(articleID uint) ([]model.Comment, error) {
	return f.listByArticleFn(articleID)
}

func TestListPublishedOnlyReturnsPublishedArticles(t *testing.T) {
	svc := service.NewArticleService(
		&fakeArticleStore{
			listFn: func(filter dao.ArticleListFilter) ([]model.Article, int64, error) {
				if !filter.PublishedOnly {
					t.Fatal("expected PublishedOnly filter to be true")
				}
				return []model.Article{
					{ID: 1, Status: service.ArticleStatusPublished},
					{ID: 2, Status: service.ArticleStatusPublished},
				}, 2, nil
			},
		},
		&fakeCategoryStore{},
		&fakeCommentReadStore{},
	)

	articles, total, _, _, err := svc.ListPublished(dto.ArticleListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected list success, got error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	for _, article := range articles {
		if article.Status != service.ArticleStatusPublished {
			t.Fatalf("expected published article, got %s", article.Status)
		}
	}
}

func TestGetPublishedFullDetailIncludesComments(t *testing.T) {
	svc := service.NewArticleService(
		&fakeArticleStore{
			findPubByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{
					ID:        id,
					Title:     "Hello",
					Status:    service.ArticleStatusPublished,
					ViewCount: 3,
					CreatedAt: time.Now(),
				}, nil
			},
			incrementFn: func(id uint) error { return nil },
		},
		&fakeCategoryStore{},
		&fakeCommentReadStore{
			listByArticleFn: func(articleID uint) ([]model.Comment, error) {
				return []model.Comment{{ID: 1, ArticleID: articleID, Content: "nice"}}, nil
			},
		},
	)

	detail, err := svc.GetPublishedFullDetail(1)
	if err != nil {
		t.Fatalf("expected full detail success, got error: %v", err)
	}
	if detail.CommentCount != 1 {
		t.Fatalf("expected comment count 1, got %d", detail.CommentCount)
	}
	if len(detail.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(detail.Comments))
	}
	if detail.Article.ViewCount != 4 {
		t.Fatalf("expected view count 4, got %d", detail.Article.ViewCount)
	}
}

func TestGetAdminDetailReturnsDraftWithoutIncreasingViewCount(t *testing.T) {
	svc := service.NewArticleService(
		&fakeArticleStore{
			findByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{ID: id, Status: service.ArticleStatusDraft, ViewCount: 3}, nil
			},
			incrementFn: func(id uint) error {
				t.Fatal("admin detail must not increment the view count")
				return nil
			},
		},
		&fakeCategoryStore{},
		&fakeCommentReadStore{},
	)

	article, err := svc.GetAdminDetail(1)
	if err != nil {
		t.Fatalf("expected admin detail success, got error: %v", err)
	}
	if article.Status != service.ArticleStatusDraft {
		t.Fatalf("expected draft article, got %s", article.Status)
	}
	if article.ViewCount != 3 {
		t.Fatalf("expected unchanged view count, got %d", article.ViewCount)
	}
}

func TestCreateArticleRejectsMissingCategory(t *testing.T) {
	svc := service.NewArticleService(
		&fakeArticleStore{
			createFn: func(article *model.Article) error { return nil },
			findByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{ID: id}, nil
			},
		},
		&fakeCategoryStore{
			findByIDFn: func(id uint) (*model.Category, error) {
				return nil, gorm.ErrRecordNotFound
			},
		},
		&fakeCommentReadStore{},
	)

	_, err := svc.Create(dto.CreateArticleRequest{
		Title:      "Test",
		Content:    "Body",
		Status:     service.ArticleStatusPublished,
		CategoryID: 999,
	}, 1)
	if err != service.ErrCategoryNotFound {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}
