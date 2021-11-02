package common

import (
	"time"

	"github.com/gtck520/ConsoleIM/models/basemodel"
	//"github.com/gtck520/ConsoleIM/models/db"

	"github.com/jinzhu/gorm"
)

//Group 用户组
type Group struct {
	basemodel.NomalModel
	Create_user  uint64 `gorm:"column:create_user;default:0;not null;" json:"create_user" form:"create_user"`              // 群主
	GroupName    string `gorm:"column:group_name;size:64;default:'';not null;" json:"group_name" form:"group_name"`        // 群名称
	Avatar       string `gorm:"column:avatar;size:255;default:'';not null;" json:"avatar" form:"avatar"`                   // 头像
	Introduction string `gorm:"column:introduction;size:255;default:'';not null;" json:"introduction" form:"introduction"` // 介绍
	MemberCount  int    `gorm:"column:user_name;size:64;default:'';not null;" json:"user_name" form:"user_name"`           // 成员数
}

// 表名
func (Group) TableName() string {
	return TableName("group")
	
}

// 添加前
func (m *Group) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *Group) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
