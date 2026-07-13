package service_test

import (
	"testing"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
)

type fakeAdminUserStore struct {
	listFn         func(filter dao.UserListFilter) ([]model.User, int64, error)
	findByIDFn     func(id uint) (*model.User, error)
	updateStatusFn func(id uint, status int8) error
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
