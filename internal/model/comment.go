package model

import "time"

type Comment struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ArticleID    uint      `gorm:"not null" json:"articleId"`
	UserID       *uint     `json:"userId"`
	GuestName    string    `gorm:"size:50" json:"guestName"`
	GuestEmail   string    `gorm:"size:255" json:"guestEmail"`
	GuestWebsite string    `gorm:"size:255" json:"guestWebsite"`
	ParentID     *uint     `json:"parentId"`
	ReplyToID    *uint     `json:"replyToId"`
	Content      string    `gorm:"size:500;not null" json:"content"`
	Status       int8      `gorm:"not null;default:1" json:"status"`
	Article      Article   `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"article"`
	User         *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Parent       *Comment  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent,omitempty"`
	ReplyTo      *Comment  `gorm:"foreignKey:ReplyToID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"replyTo,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
