package service_test

import (
	"testing"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	"blog/internal/utils"
	"gorm.io/gorm"
)

type fakeAdminUserStore struct {
	listFn           func(filter dao.UserListFilter) ([]model.User, int64, error)
	findByIDFn       func(id uint) (*model.User, error)
	findByUsernameFn func(username string) (*model.User, error)
	createFn         func(user *model.User) error
	updateStatusFn   func(id uint, status int8) error
	updateRoleFn     func(id uint, role string) error
}

func (f *fakeAdminUserStore) List(filter dao.UserListFilter) ([]model.User, int64, error) {
	return f.listFn(filter)
}

func (f *fakeAdminUserStore) FindByID(id uint) (*model.User, error) {
	return f.findByIDFn(id)
}

func (f *fakeAdminUserStore) FindByUsername(username string) (*model.User, error) {
	return f.findByUsernameFn(username)
}

func (f *fakeAdminUserStore) Create(user *model.User) error {
	return f.createFn(user)
}

func (f *fakeAdminUserStore) UpdateStatus(id uint, status int8) error {
	return f.updateStatusFn(id, status)
}

func (f *fakeAdminUserStore) UpdateRole(id uint, role string) error {
	return f.updateRoleFn(id, role)
}

func TestAdminUserListReturnsPagination(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{
		listFn: func(filter dao.UserListFilter) ([]model.User, int64, error) {
			return []model.User{{ID: 1, Username: "alice", Status: 1}}, 1, nil
		},
	})

	users, total, page, pageSize, err := svc.List(dto.AdminUserListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 1 || page != 1 || pageSize != 10 || len(users) != 1 {
		t.Fatalf("unexpected pagination result: total=%d page=%d pageSize=%d users=%d", total, page, pageSize, len(users))
	}
}

func TestUpdateUserStatusRejectsInvalidStatus(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{})

	_, err := svc.UpdateStatus(99, 1, dto.UpdateUserStatusRequest{Status: 2})
	if err != service.ErrInvalidUserStatus {
		t.Fatalf("expected ErrInvalidUserStatus, got %v", err)
	}
}

func TestUpdateUserStatusRejectsDisablingSelf(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{})
	_, err := svc.UpdateStatus(1, 1, dto.UpdateUserStatusRequest{Status: 0})
	if err != service.ErrCannotModifySelf {
		t.Fatalf("expected ErrCannotModifySelf, got %v", err)
	}
}

func TestUpdateUserRoleRejectsDemotingSelf(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{})
	_, err := svc.UpdateRole(1, 1, dto.UpdateUserRoleRequest{Role: "user"})
	if err != service.ErrCannotModifySelf {
		t.Fatalf("expected ErrCannotModifySelf, got %v", err)
	}
}

func TestUpdateUserStatusRejectsModifyingAnotherAdmin(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{
		findByIDFn: func(id uint) (*model.User, error) {
			return &model.User{ID: id, Role: "admin", Status: 1}, nil
		},
	})
	_, err := svc.UpdateStatus(1, 2, dto.UpdateUserStatusRequest{Status: 0})
	if err != service.ErrCannotModifyAdmin {
		t.Fatalf("expected ErrCannotModifyAdmin, got %v", err)
	}
}

func TestUpdateUserRoleRejectsModifyingAnotherAdmin(t *testing.T) {
	svc := service.NewUserService(&fakeAdminUserStore{
		findByIDFn: func(id uint) (*model.User, error) {
			return &model.User{ID: id, Role: "admin", Status: 1}, nil
		},
	})
	_, err := svc.UpdateRole(1, 2, dto.UpdateUserRoleRequest{Role: "user"})
	if err != service.ErrCannotModifyAdmin {
		t.Fatalf("expected ErrCannotModifyAdmin, got %v", err)
	}
}

func TestCreateAdminCreatesHashedEnabledAdmin(t *testing.T) {
	var created *model.User
	svc := service.NewUserService(&fakeAdminUserStore{
		findByUsernameFn: func(username string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(user *model.User) error {
			created = user
			return nil
		},
	})

	user, err := svc.CreateAdmin(dto.CreateAdminUserRequest{Username: "editor", Nickname: "编辑员", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if created == nil || user.Role != "admin" || user.Status != 1 {
		t.Fatalf("expected enabled admin, got %#v", user)
	}
	if user.Password == "secret123" || utils.CheckPassword(user.Password, "secret123") != nil {
		t.Fatal("expected password to be securely hashed")
	}
}
