package model

import "time"

// Tag 标签模型，用于文章的多对多关联。
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ArticleTag 文章与标签的多对多关联表。
type ArticleTag struct {
	ArticleID uint    `gorm:"primaryKey"`
	TagID     uint    `gorm:"primaryKey"`
	Article   Article `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag     `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
