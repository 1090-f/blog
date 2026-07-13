package model

import "time"

type Article struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"size:150;not null" json:"title"`
	Summary    string    `gorm:"size:255;default:''" json:"summary"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CoverImage string    `gorm:"size:255;default:''" json:"coverImage"`
	Status     string    `gorm:"size:20;not null;default:draft;index" json:"status"`
	ViewCount  int       `gorm:"not null;default:0" json:"viewCount"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	CategoryID uint      `gorm:"not null;index" json:"categoryId"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category"`
	Tags       []Tag     `gorm:"many2many:article_tags;" json:"tags"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
