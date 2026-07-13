package service_test

import (
	"errors"
	"testing"

	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	jwtpkg "blog/pkg/jwt"
	"gorm.io/gorm"
)

type fakeAuthUserStore struct {
	byUsername map[string]*model.User
	byID       map[uint]*model.User
	createFn   func(user *model.User) error
}

func (f *fakeAuthUserStore) FindByUsername(username string) (*model.User, error) {
	user, ok := f.byUsername[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeAuthUserStore) FindByID(id uint) (*model.User, error) {
	user, ok := f.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeAuthUserStore) Create(user *model.User) error {
	if f.createFn != nil {
		return f.createFn(user)
	}
	if f.byUsername == nil {
		f.byUsername = map[string]*model.User{}
	}
	if f.byID == nil {
		f.byID = map[uint]*model.User{}
	}
	user.ID = uint(len(f.byID) + 1)
	f.byUsername[user.Username] = user
	f.byID[user.ID] = user
	return nil
}

func TestRegisterHashesPassword(t *testing.T) {
	store := &fakeAuthUserStore{}
	svc := service.NewAuthService(store, "secret", 7200)

	user, err := svc.Register(dto.RegisterRequest{
		Username: "alice",
		Password: "123456",
		Nickname: "Alice",
	})
	if err != nil {
		t.Fatalf("expected register success, got error: %v", err)
	}
	if user.Password == "123456" {
		t.Fatal("expected password to be hashed")
	}
	if user.Role != "user" {
		t.Fatalf("expected role user, got %s", user.Role)
	}
}

func TestLoginReturnsToken(t *testing.T) {
	store := &fakeAuthUserStore{
		byUsername: map[string]*model.User{},
		byID:       map[uint]*model.User{},
	}
	svc := service.NewAuthService(store, "secret", 7200)
	user, err := svc.Register(dto.RegisterRequest{
		Username: "admin",
		Password: "123456",
		Nickname: "Admin",
	})
	if err != nil {
		t.Fatalf("seed register failed: %v", err)
	}
	user.Role = "admin"
	store.byUsername[user.Username] = user
	store.byID[user.ID] = user

	resp, err := svc.Login(dto.LoginRequest{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token to be returned")
	}

	claims, err := jwtpkg.ParseToken(resp.Token, "secret")
	if err != nil {
		t.Fatalf("expected token to parse, got error: %v", err)
	}
	if claims.UserID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, claims.UserID)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected role admin, got %s", claims.Role)
	}
}

func TestLoginRejectsDisabledUser(t *testing.T) {
	store := &fakeAuthUserStore{
		byUsername: map[string]*model.User{
			"disabled": {
				ID:       1,
				Username: "disabled",
				Password: "$2a$10$abc",
				Nickname: "Disabled",
				Role:     "user",
				Status:   0,
			},
		},
	}
	svc := service.NewAuthService(store, "secret", 7200)

	_, err := svc.Login(dto.LoginRequest{
		Username: "disabled",
		Password: "123456",
	})
	if !errors.Is(err, service.ErrUserDisabled) {
		t.Fatalf("expected ErrUserDisabled, got %v", err)
	}
}
