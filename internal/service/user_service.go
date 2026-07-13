package service

import (
	"errors"
	"strings"

	"blog/internal/dao"
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserStatus = errors.New("invalid user status")
	ErrInvalidProfile    = errors.New("invalid profile")
)

type UserAdminStore interface {
	List(filter dao.UserListFilter) ([]model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	UpdateStatus(id uint, status int8) error
	UpdateProfile(id uint, nickname, avatar string) error
}

type UserService struct {
	userStore UserAdminStore
}

func NewUserService(userStore UserAdminStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) List(query dto.AdminUserListQuery) ([]model.User, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)
	filter := dao.UserListFilter{
		Keyword: strings.TrimSpace(query.Keyword),
		Role:    strings.TrimSpace(query.Role),
		Status:  query.Status,
		Offset:  offset,
		Limit:   limit,
	}

	users, total, err := s.userStore.List(filter)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return users, total, page, pageSize, nil
}

func (s *UserService) UpdateStatus(id uint, req dto.UpdateUserStatusRequest) (*model.User, error) {
	if id == 0 {
		return nil, ErrUserNotFound
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, ErrInvalidUserStatus
	}

	if _, err := s.userStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if err := s.userStore.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}
	return s.userStore.FindByID(id)
}

func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*model.User, error) {
	if userID == 0 {
		return nil, ErrUserNotFound
	}

	nickname := strings.TrimSpace(req.Nickname)
	avatar := strings.TrimSpace(req.Avatar)
	if len(nickname) < 2 || len(nickname) > 50 || len(avatar) > 255 {
		return nil, ErrInvalidProfile
	}

	if _, err := s.userStore.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := s.userStore.UpdateProfile(userID, nickname, avatar); err != nil {
		return nil, err
	}
	return s.userStore.FindByID(userID)
}
