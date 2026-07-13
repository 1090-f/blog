package dao

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

type CategoryDAO struct {
	db *gorm.DB
}

func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

func (d *CategoryDAO) Create(category *model.Category) error {
	return d.db.Create(category).Error
}

func (d *CategoryDAO) List() ([]model.Category, error) {
	var categories []model.Category
	if err := d.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (d *CategoryDAO) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := d.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (d *CategoryDAO) Update(category *model.Category) error {
	return d.db.Save(category).Error
}

func (d *CategoryDAO) Delete(id uint) error {
	return d.db.Delete(&model.Category{}, id).Error
}

func (d *CategoryDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
