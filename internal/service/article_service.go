package service

import (
	"errors"
	"strings"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/utils"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
)

var (
	ErrArticleNotFound      = errors.New("article not found")
	ErrInvalidArticle       = errors.New("invalid article params")
	ErrInvalidArticleStatus = errors.New("invalid article status")
	ErrArticleInUse         = errors.New("article is in use by existing comments")
)

const (
	defaultVisitorArticleLimit = 10
	maxVisitorArticleLimit     = 20
)

type ArticleFullDetail struct {
	Article      *model.Article
	Comments     []model.Comment
	CommentCount int
}

type ArticleStore interface {
	Create(article *model.Article) error
	FindByID(id uint) (*model.Article, error)
	FindPublishedByID(id uint) (*model.Article, error)
	ListLatestPublished(limit int) ([]model.Article, error)
	ListPopularPublished(limit int) ([]model.Article, error)
	List(filter dao.ArticleListFilter) ([]model.Article, int64, error)
	Update(article *model.Article) error
	Delete(id uint) error
	CountByCategoryID(categoryID uint) (int64, error)
	IncrementViewCount(id uint) error
}

type ArticleTagStore interface {
	ListByIDs(ids []uint) ([]model.Tag, error)
	ReplaceArticleTags(articleID uint, tagIDs []uint) error
}

type ArticleCategoryStore interface {
	FindByID(id uint) (*model.Category, error)
}

type ArticleCommentStore interface {
	ListByArticleID(articleID uint) ([]model.Comment, error)
}

type ArticleService struct {
	articleDAO  ArticleStore
	categoryDAO ArticleCategoryStore
	commentDAO  ArticleCommentStore
	tagDAO      ArticleTagStore
}

func NewArticleService(articleDAO ArticleStore, categoryDAO ArticleCategoryStore, commentDAO ArticleCommentStore, tagStores ...ArticleTagStore) *ArticleService {
	var tagDAO ArticleTagStore
	if len(tagStores) > 0 {
		tagDAO = tagStores[0]
	}
	return &ArticleService{
		articleDAO:  articleDAO,
		categoryDAO: categoryDAO,
		commentDAO:  commentDAO,
		tagDAO:      tagDAO,
	}
}

func (s *ArticleService) ListPublished(query dto.ArticleListQuery) ([]model.Article, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)

	articles, total, err := s.articleDAO.List(dao.ArticleListFilter{
		CategoryID:    query.CategoryID,
		TagID:         query.TagID,
		Keyword:       strings.TrimSpace(query.Keyword),
		PublishedOnly: true,
		Offset:        offset,
		Limit:         limit,
	})
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return articles, total, page, pageSize, nil
}

func (s *ArticleService) GetPublishedDetail(id uint) (*model.Article, error) {
	return s.getPublishedDetail(id, true)
}

func (s *ArticleService) GetPublishedFullDetail(id uint) (*ArticleFullDetail, error) {
	article, err := s.getPublishedDetail(id, true)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentDAO.ListByArticleID(article.ID)
	if err != nil {
		return nil, err
	}

	return &ArticleFullDetail{
		Article:      article,
		Comments:     comments,
		CommentCount: len(comments),
	}, nil
}

func (s *ArticleService) getPublishedDetail(id uint, increaseViewCount bool) (*model.Article, error) {
	if id == 0 {
		return nil, ErrArticleNotFound
	}

	article, err := s.articleDAO.FindPublishedByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	if increaseViewCount {
		if err := s.articleDAO.IncrementViewCount(id); err != nil {
			return nil, err
		}
		article.ViewCount++
	}

	return article, nil
}

func (s *ArticleService) ListLatest(limit int) ([]model.Article, error) {
	return s.articleDAO.ListLatestPublished(normalizeVisitorArticleLimit(limit))
}

func (s *ArticleService) ListPopular(limit int) ([]model.Article, error) {
	return s.articleDAO.ListPopularPublished(normalizeVisitorArticleLimit(limit))
}

func (s *ArticleService) ListAdmin(query dto.AdminArticleListQuery) ([]model.Article, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)

	status := strings.TrimSpace(query.Status)
	if status != "" && !isValidArticleStatus(status) {
		return nil, 0, 0, 0, ErrInvalidArticleStatus
	}

	articles, total, err := s.articleDAO.List(dao.ArticleListFilter{
		CategoryID: query.CategoryID,
		Keyword:    strings.TrimSpace(query.Keyword),
		Status:     status,
		TagID:      query.TagID,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return articles, total, page, pageSize, nil
}

func (s *ArticleService) Create(req dto.CreateArticleRequest, userID uint) (*model.Article, error) {
	tagIDs, err := s.validateTagIDs(req.TagIDs)
	if err != nil {
		return nil, err
	}
	article, err := s.buildArticle(nil, req.Title, req.Summary, req.Content, req.CoverImage, req.Status, req.CategoryID, userID)
	if err != nil {
		return nil, err
	}

	if err := s.articleDAO.Create(article); err != nil {
		return nil, err
	}
	if err := s.replaceArticleTags(article.ID, tagIDs); err != nil {
		return nil, err
	}

	return s.articleDAO.FindByID(article.ID)
}

func (s *ArticleService) Update(id uint, req dto.UpdateArticleRequest, userID uint) (*model.Article, error) {
	if id == 0 {
		return nil, ErrArticleNotFound
	}

	current, err := s.articleDAO.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	var tagIDs []uint
	if req.TagIDs != nil {
		tagIDs, err = s.validateTagIDs(req.TagIDs)
		if err != nil {
			return nil, err
		}
	}

	updated, err := s.buildArticle(current, req.Title, req.Summary, req.Content, req.CoverImage, req.Status, req.CategoryID, userID)
	if err != nil {
		return nil, err
	}

	if err := s.articleDAO.Update(updated); err != nil {
		return nil, err
	}
	if req.TagIDs != nil {
		if err := s.replaceArticleTags(updated.ID, tagIDs); err != nil {
			return nil, err
		}
	}

	return s.articleDAO.FindByID(updated.ID)
}

func (s *ArticleService) validateTagIDs(tagIDs []uint) ([]uint, error) {
	for _, id := range tagIDs {
		if id == 0 {
			return nil, ErrTagNotFound
		}
	}
	normalized := NormalizeTagIDs(tagIDs)
	if len(normalized) == 0 {
		if len(tagIDs) > 0 {
			return nil, ErrTagNotFound
		}
		return normalized, nil
	}
	if s.tagDAO == nil {
		return normalized, nil
	}

	tags, err := s.tagDAO.ListByIDs(normalized)
	if err != nil {
		return nil, err
	}
	if len(tags) != len(normalized) {
		return nil, ErrTagNotFound
	}
	return normalized, nil
}

func (s *ArticleService) replaceArticleTags(articleID uint, tagIDs []uint) error {
	if s.tagDAO == nil {
		return nil
	}
	return s.tagDAO.ReplaceArticleTags(articleID, tagIDs)
}

func (s *ArticleService) Delete(id uint) error {
	if id == 0 {
		return ErrArticleNotFound
	}

	if _, err := s.articleDAO.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}

	if err := s.articleDAO.Delete(id); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			return ErrArticleInUse
		}
		return err
	}

	return nil
}

func (s *ArticleService) buildArticle(article *model.Article, title, summary, content, coverImage, status string, categoryID, userID uint) (*model.Article, error) {
	normalizedTitle := strings.TrimSpace(title)
	normalizedSummary := strings.TrimSpace(summary)
	normalizedContent := strings.TrimSpace(content)
	normalizedCoverImage := strings.TrimSpace(coverImage)
	normalizedStatus := strings.TrimSpace(status)

	if normalizedTitle == "" || normalizedContent == "" || categoryID == 0 {
		return nil, ErrInvalidArticle
	}
	if len(normalizedTitle) > 150 || len(normalizedSummary) > 255 || len(normalizedCoverImage) > 255 {
		return nil, ErrInvalidArticle
	}
	if !isValidArticleStatus(normalizedStatus) {
		return nil, ErrInvalidArticleStatus
	}
	if _, err := s.categoryDAO.FindByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if article == nil {
		if userID == 0 {
			return nil, ErrInvalidArticle
		}
		article = &model.Article{}
		article.UserID = userID
	}

	article.Title = normalizedTitle
	article.Summary = normalizedSummary
	article.Content = normalizedContent
	article.CoverImage = normalizedCoverImage
	article.Status = normalizedStatus
	article.CategoryID = categoryID

	return article, nil
}

func isValidArticleStatus(status string) bool {
	return status == ArticleStatusDraft || status == ArticleStatusPublished
}

func normalizeVisitorArticleLimit(limit int) int {
	if limit <= 0 {
		return defaultVisitorArticleLimit
	}
	if limit > maxVisitorArticleLimit {
		return maxVisitorArticleLimit
	}
	return limit
}
