package common

import (
	"time"

	"github.com/gtck520/ConsoleIM/models/basemodel"
	//"github.com/gtck520/ConsoleIM/models/db"

	"github.com/jinzhu/gorm"
)

//Group 用户组
type UserGroup struct {
	basemodel.NomalModel
	User_Id  uint64 `gorm:"column:user_id;default:0;not null;" json:"user_id" form:"user_id"`              // 用户
	Group_Id uint64 `gorm:"column:group_id;default:0;not null;" json:"group_id" form:"group_id"`           // 群
	Status   uint8  `gorm:"column:status;type:tinyint(1);default:0;not null;" json:"status" form:"status"` // 用户状态（0 申请，1 通过，2拒绝，3禁言）
}

// 表名
func (UserGroup) TableName() string {
	return TableName("user_group")

}

// 添加前
func (m *UserGroup) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *UserGroup) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
