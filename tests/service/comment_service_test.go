package service_test

import (
	"testing"

	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	"gorm.io/gorm"
)

type fakeCommentStore struct {
	createFn      func(comment *model.Comment) error
	findByIDFn    func(id uint) (*model.Comment, error)
	listByIDFn    func(articleID uint) ([]model.Comment, error)
	listAllByIDFn func(articleID uint) ([]model.Comment, error)
	deleteFn      func(id uint) error
}

func (f *fakeCommentStore) Create(comment *model.Comment) error      { return f.createFn(comment) }
func (f *fakeCommentStore) FindByID(id uint) (*model.Comment, error) { return f.findByIDFn(id) }
func (f *fakeCommentStore) ListByArticleID(articleID uint) ([]model.Comment, error) {
	return f.listByIDFn(articleID)
}
func (f *fakeCommentStore) ListAllByArticleID(articleID uint) ([]model.Comment, error) {
	return f.listAllByIDFn(articleID)
}
func (f *fakeCommentStore) Delete(id uint) error { return f.deleteFn(id) }

func TestCreateCommentRequiresPublishedArticle(t *testing.T) {
	svc := service.NewCommentService(
		&fakeCommentStore{
			createFn: func(comment *model.Comment) error { return nil },
			findByIDFn: func(id uint) (*model.Comment, error) {
				return &model.Comment{ID: id}, nil
			},
		},
		&fakeArticleStore{
			findPubByIDFn: func(id uint) (*model.Article, error) {
				return nil, gorm.ErrRecordNotFound
			},
		},
	)

	_, err := svc.Create(dto.CreateCommentRequest{
		ArticleID: 1,
		Content:   "hello",
	}, 2)
	if err != service.ErrArticleNotFound {
		t.Fatalf("expected ErrArticleNotFound, got %v", err)
	}
}

func TestCreateCommentTrimsContent(t *testing.T) {
	var created *model.Comment
	svc := service.NewCommentService(
		&fakeCommentStore{
			createFn: func(comment *model.Comment) error {
				created = comment
				comment.ID = 1
				return nil
			},
			findByIDFn: func(id uint) (*model.Comment, error) {
				return &model.Comment{ID: id, Content: created.Content}, nil
			},
		},
		&fakeArticleStore{
			findPubByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{ID: id}, nil
			},
		},
	)

	comment, err := svc.Create(dto.CreateCommentRequest{
		ArticleID: 1,
		Content:   "  hello world  ",
	}, 3)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if created.Content != "hello world" {
		t.Fatalf("expected trimmed content, got %q", created.Content)
	}
	if comment.ID != 1 {
		t.Fatalf("expected comment id 1, got %d", comment.ID)
	}
}

func TestCreateReplyCommentUsesRootParentForSecondLevel(t *testing.T) {
	var created *model.Comment
	rootID := uint(10)
	replyID := uint(11)

	svc := service.NewCommentService(
		&fakeCommentStore{
			createFn: func(comment *model.Comment) error {
				created = comment
				comment.ID = 12
				return nil
			},
			findByIDFn: func(id uint) (*model.Comment, error) {
				switch id {
				case replyID:
					return &model.Comment{ID: replyID, ArticleID: 1, ParentID: &rootID, Status: 1}, nil
				case 12:
					return &model.Comment{ID: 12, ArticleID: 1, ParentID: &rootID, ReplyToID: &replyID, Content: created.Content}, nil
				default:
					return nil, gorm.ErrRecordNotFound
				}
			},
		},
		&fakeArticleStore{
			findPubByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{ID: id}, nil
			},
		},
	)

	comment, err := svc.Create(dto.CreateCommentRequest{
		ArticleID: 1,
		ReplyToID: &replyID,
		Content:   "reply",
	}, 3)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if created.ParentID == nil || *created.ParentID != rootID {
		t.Fatalf("expected parent id %d, got %+v", rootID, created.ParentID)
	}
	if created.ReplyToID == nil || *created.ReplyToID != replyID {
		t.Fatalf("expected replyTo id %d, got %+v", replyID, created.ReplyToID)
	}
	if comment.ID != 12 {
		t.Fatalf("expected comment id 12, got %d", comment.ID)
	}
}

func TestCreateCommentRejectsReplyToSelf(t *testing.T) {
	replyID := uint(11)

	svc := service.NewCommentService(
		&fakeCommentStore{
			createFn: func(comment *model.Comment) error { return nil },
			findByIDFn: func(id uint) (*model.Comment, error) {
				return &model.Comment{ID: replyID, ArticleID: 1, UserID: 3, Status: 1}, nil
			},
		},
		&fakeArticleStore{
			findPubByIDFn: func(id uint) (*model.Article, error) {
				return &model.Article{ID: id}, nil
			},
		},
	)

	_, err := svc.Create(dto.CreateCommentRequest{
		ArticleID: 1,
		ReplyToID: &replyID,
		Content:   "self reply",
	}, 3)
	if err != service.ErrInvalidComment {
		t.Fatalf("expected ErrInvalidComment, got %v", err)
	}
}
