package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/internal/controller"
	"blog/internal/dao"
	"blog/internal/middleware"
	"blog/internal/model"
	"blog/internal/service"
	jwtpkg "blog/pkg/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type fakeAdminHTTPUserStore struct {
	users map[uint]*model.User
}

func (f *fakeAdminHTTPUserStore) FindByID(id uint) (*model.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeAdminHTTPUserStore) List(filter dao.UserListFilter) ([]model.User, int64, error) {
	list := make([]model.User, 0, len(f.users))
	for _, user := range f.users {
		list = append(list, *user)
	}
	return list, int64(len(list)), nil
}

func (f *fakeAdminHTTPUserStore) UpdateStatus(id uint, status int8) error {
	user, ok := f.users[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	user.Status = status
	return nil
}

func adminAuthToken(t *testing.T, userID uint, role string) string {
	t.Helper()
	token, err := jwtpkg.GenerateToken("secret", userID, role, 7200)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	return token
}

func TestAdminUsersRouteRejectsNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeAdminHTTPUserStore{
		users: map[uint]*model.User{
			1: {ID: 1, Username: "user", Nickname: "User", Role: "user", Status: 1},
		},
	}
	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	r := gin.New()
	admin := r.Group("/api/admin")
	admin.Use(middleware.Auth("secret", store), middleware.Admin())
	admin.GET("/users", userController.List)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminAuthToken(t, 1, "user"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestAdminUsersRouteReturnsData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := &fakeAdminHTTPUserStore{
		users: map[uint]*model.User{
			1: {ID: 1, Username: "admin", Nickname: "Admin", Role: "admin", Status: 1},
			2: {ID: 2, Username: "alice", Nickname: "Alice", Role: "user", Status: 1},
		},
	}
	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	r := gin.New()
	admin := r.Group("/api/admin")
	admin.Use(middleware.Auth("secret", store), middleware.Admin())
	admin.GET("/users", userController.List)
	admin.PUT("/users/:id/status", userController.UpdateStatus)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+adminAuthToken(t, 1, "admin"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := []byte(`{"status":0}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/admin/users/2/status", bytes.NewReader(body))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+adminAuthToken(t, 1, "admin"))
	updateRec := httptest.NewRecorder()
	r.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateRec.Code)
	}
	if store.users[2].Status != 0 {
		t.Fatalf("expected user status 0, got %d", store.users[2].Status)
	}
}
