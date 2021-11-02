package common

import (
	"time"

	"github.com/gtck520/ConsoleIM/models/basemodel"
	//"github.com/gtck520/ConsoleIM/models/db"

	"github.com/jinzhu/gorm"
)

//User 用户授权信息
type User struct {
	basemodel.NomalModel
	Invite      uint64 `gorm:"column:invite;default:0;not null;" json:"invite" form:"invite"`                               // 邀请人
	NickName    string `gorm:"column:nick_name;size:64;default:'';not null;" json:"nick_name" form:"nick_name"`             // 昵称
	Avatar      string `gorm:"column:avatar;size:255;default:'';not null;" json:"avatar" form:"avatar"`                     // 头像
	UserName    string `gorm:"column:user_name;size:64;default:'';not null;" json:"user_name" form:"user_name"`             // 用户名称
	UserPass    string `gorm:"column:user_pass;size:255;default:'';not null;" json:"user_pass" form:"user_pass"`            // 用户密码
	Phone       string `gorm:"column:phone;size:64;default:'';not null;" json:"phone" form:"phone"`                         // 电话号码
	Status      uint8  `gorm:"column:status;type:tinyint(1);default:0;not null;" json:"status" form:"status"`               // 是否禁用 （0 否，1 禁用）
	Salt        string `gorm:"column:salt;size:32;default:'';not null;" json:"salt" form:"salt"`                            // 加密盐值
	LastLogTime uint64 `gorm:"column:last_log_time;type:int;default:0;not null;" json:"last_log_time" form:"last_log_time"` // 创建时间
	LastLogIp   string `gorm:"column:last_log_ip;size:20;default:'';not null;" json:"last_log_ip" form:"last_log_ip"`       // 加密盐值
	RegIp       string `gorm:"column:reg_ip;size:20;default:'';not null;" json:"reg_ip" form:"reg_ip"`                      // 加密盐值

}

// 表名
func (User) TableName() string {
	return TableName("user")
}

// 添加前
func (m *User) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *User) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
