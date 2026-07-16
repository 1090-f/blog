package model

import "time"

// User 用户模型，存储账号、密码、角色及状态等信息。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:50;not null" json:"nickname"`
	Role      string    `gorm:"size:20;not null;default:user" json:"role"`
	Avatar    string    `gorm:"size:255;default:''" json:"avatar"`
	Status    int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
