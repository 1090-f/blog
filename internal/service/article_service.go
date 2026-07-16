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

// 文章状态常量。
const (
	ArticleStatusDraft     = "draft"     // 草稿
	ArticleStatusPublished = "published" // 已发布
)

// 文章服务相关的错误值。
var (
	ErrArticleNotFound      = errors.New("article not found")                      // 文章不存在
	ErrInvalidArticle       = errors.New("invalid article params")                 // 文章参数无效
	ErrInvalidArticleStatus = errors.New("invalid article status")                 // 文章状态无效
	ErrArticleInUse         = errors.New("article is in use by existing comments") // 文章有评论关联，无法删除
)

// 游客列表查询的默认条数和上限，防止前端传入过大值。
const (
	defaultVisitorArticleLimit = 10
	maxVisitorArticleLimit     = 20
)

// ArticleFullDetail 文章完整详情，包含评论列表和评论总数。
type ArticleFullDetail struct {
	Article      *model.Article
	Comments     []model.Comment
	CommentCount int
}

// ArticleStore 文章服务所需的持久化操作抽象。
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

// ArticleTagStore 文章服务所需的标签关联操作抽象。
type ArticleTagStore interface {
	ListByIDs(ids []uint) ([]model.Tag, error)
	ReplaceArticleTags(articleID uint, tagIDs []uint) error
}

// ArticleCategoryStore 文章服务所需的分类查询操作抽象。
type ArticleCategoryStore interface {
	FindByID(id uint) (*model.Category, error)
}

// ArticleCommentStore 文章服务所需的评论列表操作抽象。
type ArticleCommentStore interface {
	ListByArticleID(articleID uint) ([]model.Comment, error)
}

// ArticleService 文章业务逻辑层，处理文章 CRUD 及列表查询。
type ArticleService struct {
	articleDAO  ArticleStore
	categoryDAO ArticleCategoryStore
	commentDAO  ArticleCommentStore
	tagDAO      ArticleTagStore
}

// NewArticleService 创建并初始化文章业务实例。
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

// ListPublished 分页查询已发布文章列表，支持按分类、标签和关键词筛选，仅返回已发布的文章。
func (s *ArticleService) ListPublished(query dto.ArticleListQuery) ([]model.Article, int64, int, int, error) {
	// 将页码和页大小转换成数据库需要的偏移量和限制条数
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

// GetPublishedDetail 获取单篇已发布文章详情，同时递增浏览量。
func (s *ArticleService) GetPublishedDetail(id uint) (*model.Article, error) {
	return s.getPublishedDetail(id, true)
}

// GetAdminDetail 获取管理端文章详情，可返回草稿状态的文章，不递增浏览量。
func (s *ArticleService) GetAdminDetail(id uint) (*model.Article, error) {
	if id == 0 {
		return nil, ErrArticleNotFound
	}

	article, err := s.articleDAO.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return article, nil
}

// GetPublishedFullDetail 获取已发布文章的完整详情，包含评论列表和评论总数。
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

// getPublishedDetail 内部方法：获取已发布文章详情，可根据参数决定是否递增浏览量。
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

// ListLatest 查询最新的已发布文章列表，返回指定条数。
func (s *ArticleService) ListLatest(limit int) ([]model.Article, error) {
	return s.articleDAO.ListLatestPublished(normalizeVisitorArticleLimit(limit))
}

// ListPopular 查询最热门的已发布文章列表，按浏览量排序返回指定条数。
func (s *ArticleService) ListPopular(limit int) ([]model.Article, error) {
	return s.articleDAO.ListPopularPublished(normalizeVisitorArticleLimit(limit))
}

// ListAdmin 分页查询管理端文章列表，支持按状态、分类、标签和关键词筛选。
func (s *ArticleService) ListAdmin(query dto.AdminArticleListQuery) ([]model.Article, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)

	status := strings.TrimSpace(query.Status)
	// 非空状态必须是合法值，否则直接报错
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

// Create 创建新文章，校验标签、分类等参数后持久化，并关联标签。
func (s *ArticleService) Create(req dto.CreateArticleRequest, userID uint) (*model.Article, error) {
	tagIDs, err := s.validateTagIDs(req.TagIDs)
	if err != nil {
		return nil, err
	}
	article, err := s.buildArticle(nil, req.Title, req.Summary, req.Content, req.CoverImage, req.Status, req.CategoryID, userID)
	if err != nil {
		return nil, err
	}

	// 持久化文章记录
	if err := s.articleDAO.Create(article); err != nil {
		return nil, err
	}
	// 建立文章与标签的多对多关联
	if err := s.replaceArticleTags(article.ID, tagIDs); err != nil {
		return nil, err
	}

	return s.articleDAO.FindByID(article.ID)
}

// Update 更新已有文章，仅更新非空字段，若传入标签则重新关联。
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

	// 仅在请求中显式传入标签时才校验和更新
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

	// 持久化更新后的文章字段
	if err := s.articleDAO.Update(updated); err != nil {
		return nil, err
	}
	// 仅在传入标签时重新建立关联关系
	if req.TagIDs != nil {
		if err := s.replaceArticleTags(updated.ID, tagIDs); err != nil {
			return nil, err
		}
	}

	return s.articleDAO.FindByID(updated.ID)
}

// validateTagIDs 校验标签 ID 列表：去除重复、检查零值，并验证标签在数据库中真实存在。
func (s *ArticleService) validateTagIDs(tagIDs []uint) ([]uint, error) {
	// 拒绝零值 ID
	for _, id := range tagIDs {
		if id == 0 {
			return nil, ErrTagNotFound
		}
	}
	// 去重排序
	normalized := NormalizeTagIDs(tagIDs)
	if len(normalized) == 0 {
		if len(tagIDs) > 0 {
			return nil, ErrTagNotFound
		}
		return normalized, nil
	}
	// 没有注入 tagDAO 时跳过数据库校验
	if s.tagDAO == nil {
		return normalized, nil
	}

	// 查询数据库确认所有标签均存在
	tags, err := s.tagDAO.ListByIDs(normalized)
	if err != nil {
		return nil, err
	}
	if len(tags) != len(normalized) {
		return nil, ErrTagNotFound
	}
	return normalized, nil
}

// replaceArticleTags 替换文章的标签关联关系（先删旧再建新）。
func (s *ArticleService) replaceArticleTags(articleID uint, tagIDs []uint) error {
	if s.tagDAO == nil {
		return nil
	}
	return s.tagDAO.ReplaceArticleTags(articleID, tagIDs)
}

// Delete 删除指定文章，若文章存在评论关联则拒绝删除。
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
		// MySQL 错误码 1451 = 外键约束冲突，说明文章仍有评论引用
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			return ErrArticleInUse
		}
		return err
	}

	return nil
}

// buildArticle 校验文章字段并构建 Article 模型；article 为 nil 时创建新记录，否则更新已有记录。
func (s *ArticleService) buildArticle(article *model.Article, title, summary, content, coverImage, status string, categoryID, userID uint) (*model.Article, error) {
	normalizedTitle := strings.TrimSpace(title)
	normalizedSummary := strings.TrimSpace(summary)
	normalizedContent := strings.TrimSpace(content)
	normalizedCoverImage := strings.TrimSpace(coverImage)
	normalizedStatus := strings.TrimSpace(status)

	// 必填字段不能为空
	if normalizedTitle == "" || normalizedContent == "" || categoryID == 0 {
		return nil, ErrInvalidArticle
	}
	// 防止超长内容写入数据库
	if len(normalizedTitle) > 150 || len(normalizedSummary) > 255 || len(normalizedCoverImage) > 255 {
		return nil, ErrInvalidArticle
	}
	if !isValidArticleStatus(normalizedStatus) {
		return nil, ErrInvalidArticleStatus
	}
	// 校验分类是否真实存在
	if _, err := s.categoryDAO.FindByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if article == nil {
		// 新建文章时必须指定作者
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

// isValidArticleStatus 判断文章状态是否为合法值（draft 或 published）。
func isValidArticleStatus(status string) bool {
	return status == ArticleStatusDraft || status == ArticleStatusPublished
}

// normalizeVisitorArticleLimit 将游客查询的文章条数限制在合法范围内。
func normalizeVisitorArticleLimit(limit int) int {
	if limit <= 0 {
		return defaultVisitorArticleLimit
	}
	if limit > maxVisitorArticleLimit {
		return maxVisitorArticleLimit
	}
	return limit
}
