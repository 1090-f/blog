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

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentForbidden     = errors.New("comment ownership required")
	ErrInvalidComment       = errors.New("invalid comment params")
	ErrGuestNameRequired    = errors.New("guest name is required")
	ErrGuestEmailRequired   = errors.New("guest email is required")
	ErrInvalidGuestEmail    = errors.New("guest email is invalid")
	ErrInvalidGuestWebsite  = errors.New("guest website must use http or https")
	ErrInvalidCommentStatus = errors.New("invalid comment status")
)

type CommentService struct {
	commentDAO CommentStore
	articleDAO CommentArticleStore
}

type CommentStore interface {
	Create(comment *model.Comment) error
	FindByID(id uint) (*model.Comment, error)
	ListByArticleID(articleID uint) ([]model.Comment, error)
	Delete(id uint) error
}

type CommentArticleStore interface {
	FindByID(id uint) (*model.Article, error)
	FindPublishedByID(id uint) (*model.Article, error)
}

func NewCommentService(commentDAO CommentStore, articleDAO CommentArticleStore) *CommentService {
	return &CommentService{
		commentDAO: commentDAO,
		articleDAO: articleDAO,
	}
}

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

func (s *CommentService) Create(req dto.CreateCommentRequest, userID uint) (*model.Comment, error) {
	if _, err := s.articleDAO.FindPublishedByID(req.ArticleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	content := strings.TrimSpace(req.Content)
	if req.ArticleID == 0 || content == "" || len(content) > 500 {
		return nil, ErrInvalidComment
	}
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

		if replyTo.ParentID == nil {
			parentID = &replyTo.ID
		} else {
			parentID = replyTo.ParentID
		}
	}

	comment := &model.Comment{
		ArticleID: req.ArticleID,
		ParentID:  parentID,
		ReplyToID: req.ReplyToID,
		Content:   content,
		Status:    1,
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

func isValidGuestEmail(email string) bool {
	if len(email) > 255 {
		return false
	}
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

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
	if comment.UserID == nil || *comment.UserID != userID {
		return ErrCommentForbidden
	}

	return s.commentDAO.Delete(id)
}

type CommentAdminStore interface {
	ListAdmin(filter dao.CommentListFilter) ([]model.Comment, int64, error)
	FindByID(id uint) (*model.Comment, error)
	UpdateStatus(id uint, status int8) error
	Delete(id uint) error
}

type AdminCommentService struct {
	commentStore CommentAdminStore
}

func NewAdminCommentService(commentStore CommentAdminStore) *AdminCommentService {
	return &AdminCommentService{commentStore: commentStore}
}

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

func (s *AdminCommentService) UpdateStatus(id uint, req dto.UpdateCommentStatusRequest) (*model.Comment, error) {
	if id == 0 {
		return nil, ErrCommentNotFound
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, ErrInvalidCommentStatus
	}

	if _, err := s.commentStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	if err := s.commentStore.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}
	return s.commentStore.FindByID(id)
}

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
