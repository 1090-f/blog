package dao

import (
	"blog/internal/model"
	"strings"

	"gorm.io/gorm"
)

// UserListFilter 用户分页查询的筛选条件。
type UserListFilter struct {
	Keyword string
	Role    string
	Status  *int8
	Offset  int
	Limit   int
}

// UserDAO 用户表的数据访问层。
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建并初始化用户数据访问层实例。
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// FindByUsername 按用户名精确查询用户。
func (d *UserDAO) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := d.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 按主键查询用户。
func (d *UserDAO) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 新增用户记录到数据库。
func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

// List 按筛选条件分页查询用户列表，支持按关键词、角色和状态筛选，先统计总数再查询当前页数据。
func (d *UserDAO) List(filter UserListFilter) ([]model.User, int64, error) {
	query := d.db.Model(&model.User{})
	// 关键词模糊搜索（匹配用户名和昵称）
	if strings.TrimSpace(filter.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(filter.Keyword) + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", keyword, keyword)
	}
	// 按角色筛选
	if strings.TrimSpace(filter.Role) != "" {
		query = query.Where("role = ?", strings.TrimSpace(filter.Role))
	}
	// 按状态筛选（指针类型，nil 表示不筛选）
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// 先统计符合筛选条件的总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询当前页的用户数据，按创建时间倒序
	var users []model.User
	if err := query.Order("created_at DESC").Offset(filter.Offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// UpdateStatus 更新指定用户的启用/禁用状态字段。
func (d *UserDAO) UpdateStatus(id uint, status int8) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRole 更新指定用户的角色字段。
func (d *UserDAO) UpdateRole(id uint, role string) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("role", role).Error
}

// CountAll 统计所有用户的总数量。
func (d *UserDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
