package models

import (
	"time"
)

// Admin 管理员模型
type Admin struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string    `gorm:"size:255;not null" json:"-"` // 密码不返回给前端
	Email       string    `gorm:"size:100" json:"email"`
	Phone       string    `gorm:"size:20" json:"phone"`
	AvatarURL   string    `gorm:"size:255" json:"avatar_url"`
	Role        int       `gorm:"default:2;comment:1=super,2=admin,3=operator" json:"role"`
	Status      int       `gorm:"default:1;comment:1=active,2=disabled,3=deleted" json:"status"`
	Permissions string    `gorm:"type:text;comment:JSON数组" json:"permissions"`
	LastLoginAt time.Time `gorm:"comment:最后登录时间" json:"last_login_at"`
	LastLoginIP string    `gorm:"size:45;comment:最后登录IP" json:"last_login_ip"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}
