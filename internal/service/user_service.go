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

// 用户服务相关的错误哨兵值。
var (
	ErrUserNotFound      = errors.New("user not found")      // 用户不存在
	ErrInvalidUserStatus = errors.New("invalid user status") // 用户状态值无效
	ErrInvalidUserRole   = errors.New("invalid user role")   // 用户角色值无效
	ErrCannotModifySelf  = errors.New("cannot modify current administrator")  // 不能修改自己（禁用自己或降自己的权限）
	ErrCannotModifyAdmin = errors.New("administrators cannot modify other administrators") // 管理员之间不能互相修改
)

// UserAdminStore 管理端用户服务所需的持久化操作抽象。
type UserAdminStore interface {
	List(filter dao.UserListFilter) ([]model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	UpdateStatus(id uint, status int8) error
	UpdateRole(id uint, role string) error
}

// UserService 管理端用户业务逻辑层，处理用户列表和状态管理。
type UserService struct {
	userStore UserAdminStore
}

// NewUserService 创建并初始化管理端用户业务实例。
func NewUserService(userStore UserAdminStore) *UserService {
	return &UserService{userStore: userStore}
}

// List 分页查询管理端用户列表，支持按关键词、角色和状态筛选。
func (s *UserService) List(query dto.AdminUserListQuery) ([]model.User, int64, int, int, error) {
	page, pageSize, offset, limit := utils.NormalizePage(query.Page, query.PageSize)
	// 构建过滤条件，关键词和角色去除首尾空格
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

// UpdateStatus 更新用户的启用/禁用状态（0=禁用，1=启用），禁止禁用自己和管理员之间互相修改。
func (s *UserService) UpdateStatus(actorID, id uint, req dto.UpdateUserStatusRequest) (*model.User, error) {
	if id == 0 {
		return nil, ErrUserNotFound
	}
	// 状态值只能是 0（禁用）或 1（启用）
	if req.Status != 0 && req.Status != 1 {
		return nil, ErrInvalidUserStatus
	}
	// 管理员不能禁用自己
	if actorID == id && req.Status == 0 {
		return nil, ErrCannotModifySelf
	}

	// 先确认用户存在，并禁止管理员之间互相修改。
	target, err := s.userStore.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	// 管理员不能修改其他管理员
	if target.Role == "admin" && actorID != id {
		return nil, ErrCannotModifyAdmin
	}
	// 更新用户状态
	if err := s.userStore.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}
	return s.userStore.FindByID(id)
}

// CreateAdmin 创建一个启用状态的后台管理员账号。
func (s *UserService) CreateAdmin(req dto.CreateAdminUserRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	nickname := strings.TrimSpace(req.Nickname)
	password := strings.TrimSpace(req.Password)
	// 校验用户名、昵称和密码的长度范围
	if len(username) < 3 || len(username) > 50 || len(nickname) < 2 || len(nickname) > 50 || len(password) < 6 || len(password) > 32 {
		return nil, ErrInvalidRegister
	}
	// 检查用户名是否已被占用
	if _, err := s.userStore.FindByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密后存储
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// 创建启用状态的管理员账号
	user := &model.User{Username: username, Password: hashedPassword, Nickname: nickname, Role: "admin", Status: 1}
	if err := s.userStore.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateRole 调整用户角色（admin/user），禁止降低自己的权限和管理员之间互相修改。
func (s *UserService) UpdateRole(actorID, id uint, req dto.UpdateUserRoleRequest) (*model.User, error) {
	role := strings.TrimSpace(req.Role)
	// 角色只能是 admin 或 user
	if role != "admin" && role != "user" {
		return nil, ErrInvalidUserRole
	}
	// 管理员不能把自己的角色降为普通用户
	if actorID == id && role != "admin" {
		return nil, ErrCannotModifySelf
	}
	// 先确认目标用户存在
	target, err := s.userStore.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	// 管理员不能修改其他管理员的角色
	if target.Role == "admin" && actorID != id {
		return nil, ErrCannotModifyAdmin
	}
	// 更新用户角色
	if err := s.userStore.UpdateRole(id, role); err != nil {
		return nil, err
	}
	return s.userStore.FindByID(id)
}
