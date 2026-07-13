package service

import (
	"errors"
	"strings"

	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/utils"
	jwtpkg "blog/pkg/jwt"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserExists         = errors.New("username already exists")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrInvalidRegister    = errors.New("invalid register params")
)

type AuthService struct {
	userDAO       AuthUserStore
	secret        string
	expireSeconds int
}

type AuthUserStore interface {
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
}

func NewAuthService(userDAO AuthUserStore, secret string, expireSeconds int) *AuthService {
	return &AuthService{
		userDAO:       userDAO,
		secret:        secret,
		expireSeconds: expireSeconds,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	nickname := strings.TrimSpace(req.Nickname)
	password := strings.TrimSpace(req.Password)

	if len(username) < 3 || len(username) > 50 {
		return nil, ErrInvalidRegister
	}
	if len(nickname) < 2 || len(nickname) > 50 {
		return nil, ErrInvalidRegister
	}
	if len(password) < 6 || len(password) > 20 {
		return nil, ErrInvalidRegister
	}

	if _, err := s.userDAO.FindByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Nickname: nickname,
		Role:     "user",
		Status:   1,
	}

	if err := s.userDAO.Create(user); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userDAO.FindByUsername(strings.TrimSpace(req.Username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := jwtpkg.GenerateToken(s.secret, user.ID, user.Role, s.expireSeconds)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.LoginUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     user.Role,
			Avatar:   user.Avatar,
			Status:   user.Status,
		},
	}, nil
}

func (s *AuthService) Profile(userID uint) (*model.User, error) {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	return user, nil
}
