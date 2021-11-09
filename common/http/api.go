package http

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gtck520/ConsoleIM/common/logger"
	"github.com/idoubi/goz"
)

const (
	Api_Url        = "http://127.0.0.1:8000/v1/api"
	Websocket_Url  = "ws://127.0.0.1:8000/ws/chat"
	FriendList_Url = "/friend/friend_list"
	Login_Url      = "/user/login"
	UserInfo_Url   = "/user/info"
	Send_Url       = "/ws/send"
)

type Api struct {
	Token  string
	Header map[string]interface{}
	Logger logger.Logger
	Cli    *goz.Request
}
type Results struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApi() *Api {
	zap := logger.Logger{}
	zap.Init()
	header := map[string]interface{}{
		"Content-Type": "application/json",
	}
	f := &Api{
		Token:  "",
		Header: header,
		Logger: zap,
		Cli:    goz.NewClient(),
	}
	return f
}

//登录
func (a *Api) Login(username string, password string) Results {

	resp, err := a.Cli.Post(Api_Url+Login_Url, goz.Options{
		Headers: a.Header,
		JSON: struct {
			Phone    string `json:"phone"`
			Password string `json:"user_pass"`
			Code     string `json:"code"`
		}{
			username,
			password,
			"0000",
		},
	})
	if err != nil {
		a.Logger.Log.Error(err)
	}

	body, _ := resp.GetBody()
	a.Logger.Log.Info(body) //调试
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}
	result := Results{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		a.Logger.Log.Error(err)
	}
	return result
}

//用户信息
func (a *Api) Info() Results {
	resp, err := a.Cli.Post(Api_Url+UserInfo_Url, goz.Options{
		Headers: a.Header,
		JSON: struct {
		}{},
	})
	if err != nil {
		a.Logger.Log.Error(err)
	}

	body, _ := resp.GetBody()
	a.Logger.Log.Info(body) //调试
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}
	result := Results{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		a.Logger.Log.Error(err)
	}
	return result
}

//获取好友列表
func (a *Api) GetFriends() {
	resp, err := a.Cli.Post(Api_Url+FriendList_Url, goz.Options{
		Headers: a.Header,
		JSON: struct {
			Key1 string   `json:"key1"`
			Key2 []string `json:"key2"`
			Key3 int      `json:"key3"`
		}{"value1", []string{"value21", "value22"}, 333},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(body)
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}

}

//发送消息
func (a *Api) SendMessage(fromid int, toid int, typestr string, msg string) Results {
	resp, err := a.Cli.Post(Api_Url+Send_Url, goz.Options{
		Headers: a.Header,
		JSON: struct {
			FromId int    `json:"fromid"` // 发送方
			ToId   int    `json:"toid"`   // 接收方
			Type   string `json:"type"`   // 发送类型  group:群发  user:私聊
			Msg    string `json:"msg"`    // 内容
		}{fromid, toid, typestr, msg},
	})
	if err != nil {
		a.Logger.Log.Error(err)
	}

	body, _ := resp.GetBody()
	a.Logger.Log.Info(body) //调试
	result := Results{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		a.Logger.Log.Error(err)
	}
	return result

}
