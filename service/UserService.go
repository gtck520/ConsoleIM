package service

import (
	"errors"
	"time"

	"github.com/gtck520/ConsoleIM/common/codes"
	"github.com/gtck520/ConsoleIM/common/logger"
	models "github.com/gtck520/ConsoleIM/models/common"

	//pageModel "github.com/gtck520/ConsoleIM/page"
	"github.com/gtck520/ConsoleIM/common/helper"
	"github.com/gtck520/ConsoleIM/common/util/cache"
	"github.com/gtck520/ConsoleIM/common/util/convert"
	"github.com/gtck520/ConsoleIM/common/util/hash"
	"github.com/gtck520/ConsoleIM/repository"

	"github.com/gtck520/ConsoleIM/common/util/uuids"
)

// UserService
type UserService struct {
	Repository *repository.UserRepository `inject:""`
	Log        logger.ILogger             `inject:""`
}

//ExistUserByName 判断用户名是否已存在
func (u *UserService) ExistUserByPhone(phone string) bool {
	where := models.User{Phone: phone}
	isExits, _ := u.Repository.ExistUserByName(&where)
	return isExits
}

//AddUser 添加用户
func (u *UserService) AddUser(user *models.User) bool {
	//生成用户私有盐值
	user.Salt = helper.GetCode(4)
	user.UserPass = hash.Md5String(codes.MD5_PREFIX + user.Salt + user.UserPass)

	isOK := u.Repository.AddUser(user)
	if !isOK {
		return false
	}
	return true

}

//Login 用户登录
func (u *UserService) Login(user *models.User) (interface{}, error) {
	where := models.User{Phone: user.Phone}
	isExits, User := u.Repository.ExistUserByName(&where)
	if !isExits {
		return false, errors.New("用户不存在")
	} //生成用户私有盐值
	Salt := User.Salt

	user.UserPass = hash.Md5String(codes.MD5_PREFIX + Salt + user.UserPass)

	isOK, User := u.Repository.CheckUser(user)
	u.Log.Infof("%+v", User)
	if !isOK {
		return false, errors.New("用户密码错误")
	}
	// 缓存或者redis
	uuid := uuids.GetUUID()
	err := cache.Set([]byte(uuid), []byte(convert.ToString(User.ID)), 60*60) // 1H
	if err != nil {
		return false, errors.New("缓存设置失败")
	}
	// token jwt
	userInfo := make(map[string]string)
	userInfo["exp"] = convert.ToString(time.Now().Add(time.Hour * time.Duration(1)).Unix()) // 1H
	userInfo["iat"] = convert.ToString(time.Now().Unix())
	userInfo["uuid"] = uuid

	// 发至页面
	resData := make(map[string]string)
	return resData, nil

}

//GetUserById 根据id获取用户信息
func (u *UserService) GetUserById(Id uint) (interface{}, error) {
	User := u.Repository.GetUserByID(Id)
	if User == nil {
		return nil, errors.New("该用户信息不存在")
	}
	return User, nil

}
