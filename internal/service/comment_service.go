package service

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/utils"

	"gorm.io/gorm"
)

// 评论服务相关的错误值。
var (
	ErrCommentNotFound      = errors.New("comment not found")                    // 评论不存在
	ErrCommentForbidden     = errors.New("comment ownership required")           // 非评论作者，无权操作
	ErrInvalidComment       = errors.New("invalid comment params")               // 评论参数无效
	ErrGuestNameRequired    = errors.New("guest name is required")               // 游客昵称必填
	ErrGuestEmailRequired   = errors.New("guest email is required")              // 游客邮箱必填
	ErrInvalidGuestEmail    = errors.New("guest email is invalid")               // 游客邮箱格式无效
	ErrInvalidGuestWebsite  = errors.New("guest website must use http or https") // 游客网站必须为 http/https
	ErrInvalidCommentStatus = errors.New("invalid comment status")               // 评论状态值无效
)

// CommentService 公开评论业务逻辑层，处理评论查询、创建和用户删除。
type CommentService struct {
	commentDAO CommentStore
	articleDAO CommentArticleStore
}

// CommentStore 评论服务所需的持久化操作抽象。
type CommentStore interface {
	Create(comment *model.Comment) error
	FindByID(id uint) (*model.Comment, error)
	ListByArticleID(articleID uint) ([]model.Comment, error)
	Delete(id uint) error
}

// CommentArticleStore 评论服务所需的文章查询操作抽象。
type CommentArticleStore interface {
	FindByID(id uint) (*model.Article, error)
	FindPublishedByID(id uint) (*model.Article, error)
}

// NewCommentService 创建并初始化公开评论业务实例。
func NewCommentService(commentDAO CommentStore, articleDAO CommentArticleStore) *CommentService {
	return &CommentService{
		commentDAO: commentDAO,
		articleDAO: articleDAO,
	}
}

// ListPublicByArticleID 查询指定已发布文章的公开评论列表，先校验文章是否存在且已发布。
func (s *CommentService) ListPublicByArticleID(articleID uint) ([]model.Comment, error) {
	if articleID == 0 {
		return nil, ErrArticleNotFound
	}

	if _, err := s.articleDAO.FindPublishedByID(articleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return s.commentDAO.ListByArticleID(articleID)
}

// Create 在已发布文章下创建评论，支持游客和登录用户两种身份；若为回复评论则自动关联父级。
func (s *CommentService) Create(req dto.CreateCommentRequest, userID uint) (*model.Comment, error) {
	if _, err := s.articleDAO.FindPublishedByID(req.ArticleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	content := strings.TrimSpace(req.Content)
	// 文章 ID 不能为空，评论内容不能为空且不能超过 500 字
	if req.ArticleID == 0 || content == "" || len(content) > 500 {
		return nil, ErrInvalidComment
	}
	// 游客身份必须填写昵称和邮箱，邮箱格式和网站地址需校验
	guestName := strings.TrimSpace(req.GuestName)
	guestEmail := strings.TrimSpace(req.GuestEmail)
	guestWebsite := strings.TrimSpace(req.GuestWebsite)
	if userID == 0 {
		if guestName == "" || len(guestName) > 50 {
			return nil, ErrGuestNameRequired
		}
		if guestEmail == "" {
			return nil, ErrGuestEmailRequired
		}
		if !isValidGuestEmail(guestEmail) {
			return nil, ErrInvalidGuestEmail
		}
		if !isValidGuestWebsite(guestWebsite) {
			return nil, ErrInvalidGuestWebsite
		}
	}

	var parentID *uint
	if req.ReplyToID != nil {
		// 回复评论：校验被回复评论存在、属于同一篇文章、状态正常、且不能回复自己
		replyTo, err := s.commentDAO.FindByID(*req.ReplyToID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCommentNotFound
			}
			return nil, err
		}
		if replyTo.ArticleID != req.ArticleID || replyTo.Status != 1 {
			return nil, ErrInvalidComment
		}
		if userID != 0 && replyTo.UserID != nil && *replyTo.UserID == userID {
			return nil, ErrInvalidComment
		}

		// 将二级回复统一挂到根评论下，保持最多两层嵌套
		if replyTo.ParentID == nil {
			parentID = &replyTo.ID
		} else {
			parentID = replyTo.ParentID
		}
	}

	// 构建评论模型，登录用户和游客填写不同字段
	comment := &model.Comment{
		ArticleID: req.ArticleID,
		ParentID:  parentID,
		ReplyToID: req.ReplyToID,
		Content:   content,
		Status:    1, // 默认状态为待审核
	}
	if userID != 0 {
		comment.UserID = &userID
	} else {
		comment.GuestName = guestName
		comment.GuestEmail = guestEmail
		comment.GuestWebsite = guestWebsite
	}
	if err := s.commentDAO.Create(comment); err != nil {
		return nil, err
	}

	return s.commentDAO.FindByID(comment.ID)
}

// isValidGuestEmail 校验游客邮箱格式是否合法（长度不超过 255 且符合 RFC 5322）。
func isValidGuestEmail(email string) bool {
	if len(email) > 255 {
		return false
	}
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

// isValidGuestWebsite 校验游客网站地址是否合法（空值允许，非空时必须为 http/https 协议）。
func isValidGuestWebsite(website string) bool {
	if website == "" {
		return true
	}
	if len(website) > 255 {
		return false
	}

	parsed, err := url.Parse(website)
	if err != nil || parsed.Host == "" {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return scheme == "http" || scheme == "https"
}

// DeleteByUser 用户删除自己的评论，校验评论归属后执行删除。
func (s *CommentService) DeleteByUser(id, userID uint) error {
	if id == 0 || userID == 0 {
		return ErrCommentNotFound
	}

	comment, err := s.commentDAO.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	// 只有评论作者本人可以删除自己的评论
	if comment.UserID == nil || *comment.UserID != userID {
		return ErrCommentForbidden
	}

	return s.commentDAO.Delete(id)
}

// CommentAdminStore 管理端评论服务所需的持久化操作抽象。
type CommentAdminStore interface {
	ListAdmin(filter dao.CommentListFilter) ([]model.Comment, int64, error)
	FindByID(id uint) (*model.Comment, error)
	UpdateStatus(id uint, status int8) error
	Delete(id uint) error
}

// AdminCommentService 管理端评论业务逻辑层，处理评论审核和删除。
type AdminCommentService struct {
	commentStore CommentAdminStore
}

// NewAdminCommentService 创建并初始化管理端评论业务实例。
func NewAdminCommentService(commentStore CommentAdminStore) *AdminCommentService {
	return &AdminCommentService{commentStore: commentStore}
}

// List 分页查询管理端评论列表，支持按关键词、文章 ID 和状态筛选。
func (s *AdminCommentService) List(query dto.AdminCommentListQuery) ([]model.Comment, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)
	comments, total, err := s.commentStore.ListAdmin(dao.CommentListFilter{
		Keyword:   strings.TrimSpace(query.Keyword),
		ArticleID: query.ArticleID,
		Status:    query.Status,
		Offset:    offset,
		Limit:     limit,
	})
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return comments, total, page, pageSize, nil
}

// UpdateStatus 更新评论的审核状态（0=隐藏，1=显示），校验参数后持久化。
func (s *AdminCommentService) UpdateStatus(id uint, req dto.UpdateCommentStatusRequest) (*model.Comment, error) {
	if id == 0 {
		return nil, ErrCommentNotFound
	}
	// 状态值只能是 0（隐藏）或 1（显示）
	if req.Status != 0 && req.Status != 1 {
		return nil, ErrInvalidCommentStatus
	}

	// 先确认评论存在
	if _, err := s.commentStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	// 更新审核状态
	if err := s.commentStore.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}
	return s.commentStore.FindByID(id)
}

// Delete 管理端删除指定评论，校验评论存在后执行删除。
func (s *AdminCommentService) Delete(id uint) error {
	if id == 0 {
		return ErrCommentNotFound
	}
	if _, err := s.commentStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	return s.commentStore.Delete(id)
}
