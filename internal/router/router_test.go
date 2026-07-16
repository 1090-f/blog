package router

import (
	"testing"

	"blog/config"
	"github.com/gin-gonic/gin"
)

func TestPublicServerDoesNotRegisterAdminRoutes(t *testing.T) {
	routes := routeSet(New(config.Config{}, nil))

	if !routes["POST /api/auth/register"] {
		t.Fatal("public server must register user registration")
	}
	if !routes["GET /api/articles"] {
		t.Fatal("public server must register public article routes")
	}
	if routes["GET /api/admin/dashboard"] || routes["POST /api/admin/login"] {
		t.Fatal("public server must not register administrator routes")
	}
}

func TestAdminServerDoesNotRegisterPublicRoutes(t *testing.T) {
	routes := routeSet(NewAdmin(config.Config{}, nil))

	for _, route := range []string{
		"POST /api/admin/login",
		"GET /api/admin/session",
		"POST /api/admin/users/admin",
		"PUT /api/admin/users/:id/role",
		"GET /api/admin/articles/:id",
		"GET /api/admin/categories",
		"GET /api/admin/tags",
	} {
		if !routes[route] {
			t.Fatalf("admin server must register %s", route)
		}
	}

	for _, route := range []string{
		"POST /api/auth/register",
		"POST /api/auth/login",
		"GET /api/articles",
		"POST /api/comments",
		"GET /api/user/session",
	} {
		if routes[route] {
			t.Fatalf("admin server must not register %s", route)
		}
	}
}

func routeSet(engine *gin.Engine) map[string]bool {
	routes := make(map[string]bool)
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	return routes
}
