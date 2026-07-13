package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"blog/internal/controller"
	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/middleware"
	"blog/internal/model"
	"blog/internal/service"
	jwtpkg "blog/pkg/jwt"
	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type fakeHTTPCommentStore struct {
	comments map[uint]*model.Comment
	nextID   uint
}

func (f *fakeHTTPCommentStore) Create(comment *model.Comment) error {
	if f.comments == nil {
		f.comments = map[uint]*model.Comment{}
	}
	f.nextID++
	comment.ID = f.nextID
	f.comments[comment.ID] = comment
	return nil
}

func (f *fakeHTTPCommentStore) FindByID(id uint) (*model.Comment, error) {
	comment, ok := f.comments[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return comment, nil
}

func (f *fakeHTTPCommentStore) ListByArticleID(articleID uint) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	for _, comment := range f.comments {
		if comment.ArticleID == articleID && comment.Status == 1 {
			comments = append(comments, *comment)
		}
	}
	return comments, nil
}

func (f *fakeHTTPCommentStore) Delete(id uint) error {
	if _, ok := f.comments[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(f.comments, id)
	return nil
}

func (f *fakeHTTPCommentStore) ListAdmin(filter dao.CommentListFilter) ([]model.Comment, int64, error) {
	comments := make([]model.Comment, 0, len(f.comments))
	for _, comment := range f.comments {
		if filter.ArticleID != 0 && comment.ArticleID != filter.ArticleID {
			continue
		}
		if filter.Status != nil && comment.Status != *filter.Status {
			continue
		}
		comments = append(comments, *comment)
	}
	sort.Slice(comments, func(i, j int) bool { return comments[i].ID > comments[j].ID })
	return comments, int64(len(comments)), nil
}

func (f *fakeHTTPCommentStore) UpdateStatus(id uint, status int8) error {
	comment, ok := f.comments[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	comment.Status = status
	return nil
}

type fakeHTTPCommentArticleStore struct {
	articles map[uint]*model.Article
}

func (f *fakeHTTPCommentArticleStore) FindByID(id uint) (*model.Article, error) {
	article, ok := f.articles[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return article, nil
}

func (f *fakeHTTPCommentArticleStore) FindPublishedByID(id uint) (*model.Article, error) {
	return f.FindByID(id)
}

type fakeHTTPCommentUserReader struct {
	users map[uint]*model.User
}

func (f *fakeHTTPCommentUserReader) FindByID(id uint) (*model.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func newCommentHTTPRouter(store *fakeHTTPCommentStore, articles *fakeHTTPCommentArticleStore, users *fakeHTTPCommentUserReader) *gin.Engine {
	commentController := controller.NewCommentController(service.NewCommentService(store, articles))
	adminCommentController := controller.NewAdminCommentController(service.NewAdminCommentService(store))

	r := gin.New()
	api := r.Group("/api")
	api.POST("/comments", middleware.OptionalAuth("secret", users), commentController.Create)

	admin := api.Group("/admin")
	admin.Use(middleware.Auth("secret", users), middleware.Admin())
	admin.GET("/comments", adminCommentController.List)
	return r
}

func newCommentHTTPRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/comments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestGuestCommentRouteCreatesComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeHTTPCommentStore{}
	r := newCommentHTTPRouter(store, &fakeHTTPCommentArticleStore{articles: map[uint]*model.Article{1: {ID: 1}}}, &fakeHTTPCommentUserReader{})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, newCommentHTTPRequest(`{"articleId":1,"content":" Nice post ","guestName":" Visitor ","guestEmail":" visitor@example.com ","guestWebsite":" https://example.com "}`))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		Code int                 `json:"code"`
		Data dto.CommentResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if response.Code != 0 || response.Data.UserID != nil || response.Data.Author.Nickname != "Visitor" || response.Data.Author.Website != "https://example.com" {
		t.Fatalf("unexpected guest response: %+v", response)
	}
	if created := store.comments[1]; created == nil || created.GuestEmail != "visitor@example.com" || created.Content != "Nice post" {
		t.Fatalf("guest comment was not persisted correctly: %+v", created)
	}
}

func TestGuestCommentRouteRequiresNameAndEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeHTTPCommentStore{}
	r := newCommentHTTPRouter(store, &fakeHTTPCommentArticleStore{articles: map[uint]*model.Article{1: {ID: 1}}}, &fakeHTTPCommentUserReader{})

	tests := []struct {
		name string
		body string
	}{
		{name: "missing name", body: `{"articleId":1,"content":"hello","guestEmail":"visitor@example.com"}`},
		{name: "missing email", body: `{"articleId":1,"content":"hello","guestName":"Visitor"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, newCommentHTTPRequest(tt.body))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGuestCommentRouteRejectsInvalidEmailAndWebsite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeHTTPCommentStore{}
	r := newCommentHTTPRouter(store, &fakeHTTPCommentArticleStore{articles: map[uint]*model.Article{1: {ID: 1}}}, &fakeHTTPCommentUserReader{})

	tests := []struct {
		name string
		body string
	}{
		{name: "invalid email", body: `{"articleId":1,"content":"hello","guestName":"Visitor","guestEmail":"visitor@"}`},
		{name: "unsupported website scheme", body: `{"articleId":1,"content":"hello","guestName":"Visitor","guestEmail":"visitor@example.com","guestWebsite":"javascript:alert(1)"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, newCommentHTTPRequest(tt.body))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGuestCommentRouteRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeHTTPCommentStore{}
	r := newCommentHTTPRouter(store, &fakeHTTPCommentArticleStore{articles: map[uint]*model.Article{1: {ID: 1}}}, &fakeHTTPCommentUserReader{})

	expiredToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtpkg.Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	})
	token, err := expiredToken.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}

	req := newCommentHTTPRequest(`{"articleId":1,"content":"hello","guestName":"Visitor","guestEmail":"visitor@example.com"}`)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAdminCommentListRouteIncludesGuestComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeHTTPCommentStore{comments: map[uint]*model.Comment{
		1: {
			ID:        1,
			ArticleID: 1,
			Article:   model.Article{ID: 1, Title: "Guest post"},
			GuestName: "Visitor",
			Status:    1,
			Content:   "Hello",
		},
	}, nextID: 1}
	users := &fakeHTTPCommentUserReader{users: map[uint]*model.User{
		1: {ID: 1, Username: "admin", Nickname: "Admin", Role: "admin", Status: 1},
	}}
	r := newCommentHTTPRouter(store, &fakeHTTPCommentArticleStore{}, users)

	token, err := jwtpkg.GenerateToken("secret", 1, "admin", 3600)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/admin/comments", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		Data dto.AdminCommentListResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(response.Data.List) != 1 || response.Data.List[0].Nickname != "Visitor" || response.Data.List[0].Username != "" || response.Data.List[0].ArticleTitle != "Guest post" {
		t.Fatalf("unexpected admin comment response: %+v", response.Data.List)
	}
}
