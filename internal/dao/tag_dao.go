package dao

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

// TagDAO 标签表及标签关联表的数据访问层。
type TagDAO struct {
	db *gorm.DB
}

// NewTagDAO 创建并初始化标签数据访问层实例。
func NewTagDAO(db *gorm.DB) *TagDAO {
	return &TagDAO{db: db}
}

// List 查询所有标签，按名称升序排列。
func (d *TagDAO) List() ([]model.Tag, error) {
	var tags []model.Tag
	if err := d.db.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CountAll 统计所有标签的总数量。
func (d *TagDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Tag{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListByIDs 按主键列表批量查询标签，空列表时直接返回空结果。
func (d *TagDAO) ListByIDs(ids []uint) ([]model.Tag, error) {
	if len(ids) == 0 {
		return []model.Tag{}, nil
	}

	var tags []model.Tag
	if err := d.db.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Create 新增标签记录到数据库。
func (d *TagDAO) Create(tag *model.Tag) error {
	return d.db.Create(tag).Error
}

// FindByID 按主键查询标签。
func (d *TagDAO) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := d.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update 更新标签记录（全量保存）。
func (d *TagDAO) Update(tag *model.Tag) error {
	return d.db.Save(tag).Error
}

// Delete 按主键删除标签记录。
func (d *TagDAO) Delete(id uint) error {
	return d.db.Delete(&model.Tag{}, id).Error
}

// CountByTagID 统计指定标签在 article_tags 关联表中的引用次数。
func (d *TagDAO) CountByTagID(id uint) (int64, error) {
	var count int64
	if err := d.db.Model(&model.ArticleTag{}).Where("tag_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ReplaceArticleTags 在事务中替换文章的标签关联：先删除旧关联，再批量插入新关联。
func (d *TagDAO) ReplaceArticleTags(articleID uint, tagIDs []uint) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 删除该文章的所有旧标签关联
		if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		// 新标签列表为空时直接返回
		if len(tagIDs) == 0 {
			return nil
		}

		// 批量插入新的标签关联记录
		rows := make([]model.ArticleTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			rows = append(rows, model.ArticleTag{ArticleID: articleID, TagID: tagID})
		}
		return tx.Create(&rows).Error
	})
}
