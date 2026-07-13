package dao

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

type TagDAO struct {
	db *gorm.DB
}

func NewTagDAO(db *gorm.DB) *TagDAO {
	return &TagDAO{db: db}
}

func (d *TagDAO) List() ([]model.Tag, error) {
	var tags []model.Tag
	if err := d.db.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (d *TagDAO) CountAll() (int64, error) {
	var count int64
	if err := d.db.Model(&model.Tag{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

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

func (d *TagDAO) Create(tag *model.Tag) error {
	return d.db.Create(tag).Error
}

func (d *TagDAO) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := d.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (d *TagDAO) Update(tag *model.Tag) error {
	return d.db.Save(tag).Error
}

func (d *TagDAO) Delete(id uint) error {
	return d.db.Delete(&model.Tag{}, id).Error
}

func (d *TagDAO) CountByTagID(id uint) (int64, error) {
	var count int64
	if err := d.db.Model(&model.ArticleTag{}).Where("tag_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *TagDAO) ReplaceArticleTags(articleID uint, tagIDs []uint) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		if len(tagIDs) == 0 {
			return nil
		}

		rows := make([]model.ArticleTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			rows = append(rows, model.ArticleTag{ArticleID: articleID, TagID: tagID})
		}
		return tx.Create(&rows).Error
	})
}
