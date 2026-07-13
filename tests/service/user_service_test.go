package service_test

import (
	"testing"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
)

type fakeAdminUserStore struct {
	listFn          func(filter dao.UserListFilter) ([]model.User, int64, error)
	findByIDFn      func(id uint) (*model.User, error)
	updateStatusFn  func(id uint, status int8) error
	updateProfileFn func(id uint, nickname, avatar string) error
}

func (f *fakeAdminUserStore) List(filter dao.UserListFilter) ([]model.User, int64, error) {
	return f.listFn(filter)
}

func (f *fakeAdminUserStore) FindByID(id uint) (*model.User, error) {
	return f.findByIDFn(id)
}

func (f *fakeAdminUserStore) UpdateStatus(id uint, status int8) error {
	return f.updateStatusFn(id, status)
}

func (f *fakeAdminUserStore) UpdateProfile(id uint, nickname, avatar string) error {
	return f.updateProfileFn(id, nickname, avatar)
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

	_, err := svc.UpdateStatus(1, dto.UpdateUserStatusRequest{Status: 2})
	if err != service.ErrInvalidUserStatus {
		t.Fatalf("expected ErrInvalidUserStatus, got %v", err)
	}
}

func TestUpdateProfileTrimsAndPersistsFields(t *testing.T) {
	store := &fakeAdminUserStore{
		findByIDFn: func(id uint) (*model.User, error) {
			return &model.User{ID: id, Username: "alice", Nickname: "Alice", Avatar: ""}, nil
		},
		updateProfileFn: func(id uint, nickname, avatar string) error {
			if id != 1 {
				t.Fatalf("expected id 1, got %d", id)
			}
			if nickname != "Alice Cooper" {
				t.Fatalf("expected trimmed nickname, got %q", nickname)
			}
			if avatar != "/uploads/avatar.png" {
				t.Fatalf("expected avatar persisted, got %q", avatar)
			}
			return nil
		},
	}

	svc := service.NewUserService(store)
	user, err := svc.UpdateProfile(1, dto.UpdateProfileRequest{
		Nickname: "  Alice Cooper  ",
		Avatar:   " /uploads/avatar.png ",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("expected user id 1, got %d", user.ID)
	}
}
