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

// 认证服务相关的错误值。
var (
	ErrInvalidCredentials = errors.New("invalid username or password") // 用户名或密码错误
	ErrUserExists         = errors.New("username already exists")      // 用户名已存在
	ErrUserDisabled       = errors.New("user is disabled")             // 用户已被禁用
	ErrInvalidRegister    = errors.New("invalid register params")      // 注册参数无效
)

// AuthService 认证业务逻辑层，处理注册、登录和会话查询。
type AuthService struct {
	userDAO       AuthUserStore
	secret        string
	expireSeconds int
}

// AuthUserStore 认证服务所需的用户持久化操作抽象。
type AuthUserStore interface {
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
}

// NewAuthService 创建并初始化认证业务实例。
func NewAuthService(userDAO AuthUserStore, secret string, expireSeconds int) *AuthService {
	return &AuthService{
		userDAO:       userDAO,
		secret:        secret,
		expireSeconds: expireSeconds,
	}
}

// Register 注册用户。
func (s *AuthService) Register(req dto.RegisterRequest) (*model.User, error) {
	// 去除首尾空格，防止用户误输入空格
	username := strings.TrimSpace(req.Username)
	nickname := strings.TrimSpace(req.Nickname)
	password := strings.TrimSpace(req.Password)

	// 参数校验
	if len(username) < 3 || len(username) > 50 {
		return nil, ErrInvalidRegister
	}
	if len(nickname) < 2 || len(nickname) > 50 {
		return nil, ErrInvalidRegister
	}
	if len(password) < 6 || len(password) > 20 {
		return nil, ErrInvalidRegister
	}

	// 检查用户名是否存在
	if _, err := s.userDAO.FindByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密
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

// Login 验证用户凭据并签发登录令牌。
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查询用户
	user, err := s.userDAO.FindByUsername(strings.TrimSpace(req.Username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 检查账号状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 验证密码
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成 jwt token
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

// CurrentUser 查询当前登录用户。
func (s *AuthService) CurrentUser(userID uint) (*model.User, error) {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	return user, nil
}
