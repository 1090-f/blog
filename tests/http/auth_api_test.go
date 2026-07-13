package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/internal/controller"
	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/middleware"
	"blog/internal/model"
	"blog/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type fakeHTTPAuthUserStore struct {
	usersByUsername map[string]*model.User
	usersByID       map[uint]*model.User
}

func (f *fakeHTTPAuthUserStore) FindByUsername(username string) (*model.User, error) {
	user, ok := f.usersByUsername[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeHTTPAuthUserStore) FindByID(id uint) (*model.User, error) {
	user, ok := f.usersByID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeHTTPAuthUserStore) List(filter dao.UserListFilter) ([]model.User, int64, error) {
	list := make([]model.User, 0, len(f.usersByID))
	for _, user := range f.usersByID {
		list = append(list, *user)
	}
	return list, int64(len(list)), nil
}

func (f *fakeHTTPAuthUserStore) Create(user *model.User) error {
	if f.usersByUsername == nil {
		f.usersByUsername = map[string]*model.User{}
	}
	if f.usersByID == nil {
		f.usersByID = map[uint]*model.User{}
	}
	user.ID = uint(len(f.usersByID) + 1)
	f.usersByUsername[user.Username] = user
	f.usersByID[user.ID] = user
	return nil
}

func (f *fakeHTTPAuthUserStore) UpdateStatus(id uint, status int8) error {
	user, ok := f.usersByID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	user.Status = status
	return nil
}

func (f *fakeHTTPAuthUserStore) UpdateProfile(id uint, nickname, avatar string) error {
	user, ok := f.usersByID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	user.Nickname = nickname
	user.Avatar = avatar
	return nil
}

func TestRegisterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeHTTPAuthUserStore{}
	authService := service.NewAuthService(store, "secret", 7200)
	authController := controller.NewAuthController(authService)

	r := gin.New()
	r.POST("/api/auth/register", authController.Register)

	body := []byte(`{"username":"alice","password":"123456","nickname":"Alice"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestProfileRouteRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeHTTPAuthUserStore{}
	authService := service.NewAuthService(store, "secret", 7200)
	authController := controller.NewAuthController(authService)

	r := gin.New()
	protected := r.Group("/api/user")
	protected.Use(middleware.Auth("secret", store))
	protected.GET("/profile", authController.Profile)

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestLoginThenProfileRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeHTTPAuthUserStore{}
	authService := service.NewAuthService(store, "secret", 7200)
	authController := controller.NewAuthController(authService)

	registerReq := dto.RegisterRequest{Username: "alice", Password: "123456", Nickname: "Alice"}
	user, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("seed register failed: %v", err)
	}
	store.usersByID[user.ID] = user

	r := gin.New()
	r.POST("/api/auth/login", authController.Login)
	protected := r.Group("/api/user")
	protected.Use(middleware.Auth("secret", store))
	protected.GET("/profile", authController.Profile)

	loginBody := []byte(`{"username":"alice","password":"123456"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	r.ServeHTTP(loginRec, loginReq)

	var loginResp struct {
		Data struct {
			Token string `json:"token"`
			User  struct {
				Username string `json:"username"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response failed: %v", err)
	}
	if loginResp.Data.User.Username != "alice" {
		t.Fatalf("expected login response user alice, got %q", loginResp.Data.User.Username)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Data.Token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateProfileRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeHTTPAuthUserStore{}
	authService := service.NewAuthService(store, "secret", 7200)
	authController := controller.NewAuthController(authService)
	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	registerReq := dto.RegisterRequest{Username: "alice", Password: "123456", Nickname: "Alice"}
	user, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("seed register failed: %v", err)
	}
	store.usersByID[user.ID] = user

	r := gin.New()
	r.POST("/api/auth/login", authController.Login)
	protected := r.Group("/api/user")
	protected.Use(middleware.Auth("secret", store))
	protected.PUT("/profile", userController.UpdateProfile)

	loginBody := []byte(`{"username":"alice","password":"123456"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	r.ServeHTTP(loginRec, loginReq)

	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response failed: %v", err)
	}

	updateBody := []byte(`{"nickname":"Alice New","avatar":"/uploads/avatar.png"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Data.Token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	updated, err := store.FindByID(user.ID)
	if err != nil {
		t.Fatalf("find updated user failed: %v", err)
	}
	if updated.Nickname != "Alice New" || updated.Avatar != "/uploads/avatar.png" {
		t.Fatalf("profile not updated: nickname=%q avatar=%q", updated.Nickname, updated.Avatar)
	}
}
