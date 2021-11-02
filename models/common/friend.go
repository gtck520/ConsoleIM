package common

import (
	"time"

	"github.com/gtck520/ConsoleIM/models/basemodel"
	//"github.com/gtck520/ConsoleIM/models/db"

	"github.com/jinzhu/gorm"
)

//friend 朋友关系
type Friend struct {
	basemodel.NomalModel
	UserId   uint64    `gorm:"column:user_id;default:0;not null;" json:"user_id" form:"user_id"`              // 好友发起人
	FriendId uint64    `gorm:"column:friend_id;default:0;not null;" json:"friend_id" form:"friend_id"`        // 被请求好友
	Status   uint8     `gorm:"column:status;type:tinyint(1);default:0;not null;" json:"status" form:"status"` // 好友状态（0 申请，1 通过，2拒绝，3单向拉黑，4双向拉黑）
	FriendOn time.Time `gorm:"column:friend_on;default:0;not null;" json:"friend_on" form:"friend_on"`        //成为好友时间
}

// 表名
func (Friend) TableName() string {
	return TableName("friend")
}

// 添加前
func (m *Friend) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *Friend) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
