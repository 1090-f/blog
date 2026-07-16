package dao

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

// CategoryDAO 分类表的数据访问层。
type CategoryDAO struct {
	db *gorm.DB
}

// NewCategoryDAO 创建并初始化分类数据实例。
func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

// Create 创建分类数据记录。
func (d *CategoryDAO) Create(category *model.Category) error {
	return d.db.Create(category).Error
}

// List 查询分类数据列表。
func (d *CategoryDAO) List() ([]model.Category, error) {
	var categories []model.Category
	if err := d.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID 按条件查询分类数据。
func (d *CategoryDAO) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := d.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Update 更新分类数据记录。
func (d *CategoryDAO) Update(category *model.Category) error {
	return d.db.Save(category).Error
}

// Delete 删除分类数据记录。
func (d *CategoryDAO) Delete(id uint) error {
	return d.db.Delete(&model.Category{}, id).Error
}

// CountAll 统计分类数据数量。
func (d *CategoryDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
