package model

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ArticleTag struct {
	ArticleID uint    `gorm:"primaryKey"`
	TagID     uint    `gorm:"primaryKey"`
	Article   Article `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag     `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
