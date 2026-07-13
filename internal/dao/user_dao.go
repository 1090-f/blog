package dao

import (
	"blog/internal/model"
	"strings"

	"gorm.io/gorm"
)

type UserListFilter struct {
	Keyword string
	Role    string
	Status  *int8
	Offset  int
	Limit   int
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := d.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) List(filter UserListFilter) ([]model.User, int64, error) {
	query := d.db.Model(&model.User{})
	if strings.TrimSpace(filter.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(filter.Keyword) + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", keyword, keyword)
	}
	if strings.TrimSpace(filter.Role) != "" {
		query = query.Where("role = ?", strings.TrimSpace(filter.Role))
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	if err := query.Order("created_at DESC").Offset(filter.Offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (d *UserDAO) UpdateStatus(id uint, status int8) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (d *UserDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
