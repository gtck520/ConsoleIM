package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gtck520/ConsoleIM/common/codes"
	"github.com/gtck520/ConsoleIM/common/logger"
	"github.com/gtck520/ConsoleIM/common/util/convert"
	service "github.com/gtck520/ConsoleIM/service/v1/api"

	"log"
	"net/http"
)

//服务管理
type ServerManager struct {
	Connlist      map[string]*User //所有连接通道
	register      chan *User       //注册客户端消息通道
	ClientMessage chan *Message    //发送客户端消息通道
	ServerMessage chan *Message    //发送服务端消息通道
	SysBroadcast  chan *Message    //系统广播通道
	Len           int              //长度
	Lock          sync.Mutex
	UserService   *service.UserService `inject:""`
	Log           logger.ILogger       `inject:""`
}
type User struct {
	Conn     *websocket.Conn
	Failtime int //发送错误次数
	Name     string
	Id       string
	Avator   string
	To_id    string
	group_id string
}

// 消息载体
type Message struct {
	user        *User
	context     *gin.Context
	content     []byte
	messageType int //消息类型 文本、图片、等等
}

//发送内容content 载体
type TypeMessage struct {
	Type interface{} `json:"type"` //内容分发类型：ping、init、message 等等
	Data interface{} `json:"data"`
}

//TypeMessage Data
type ClientMessage struct {
	Name      string `json:"name"`
	Avator    string `json:"avator"`
	Id        string `json:"id"`
	VisitorId string `json:"visitor_id"`
	Group     string `json:"group"`
	Time      string `json:"time"`
	ToId      string `json:"to_id"`
	Content   string `json:"content"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	Refer     string `json:"refer"`
}

//客户端连接列表
var upgrader = websocket.Upgrader{}
var SocketServer ServerManager

func init() {
	//初始化socket参数
	SocketServer = ServerManager{
		Connlist:      make(map[string]*User),
		register:      make(chan *User, 128),
		ClientMessage: make(chan *Message, 128),
		ServerMessage: make(chan *Message, 128),
		SysBroadcast:  make(chan *Message, 128),
		Len:           0,
	}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
func (s *ServerManager) NewChatServer(c *gin.Context) {
	UserId, ok := c.Get(codes.USER_ID_Key)
	log.Println("userid:", UserId)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//客户端连接成功
	log.Print("客户端连接成功：" + conn.RemoteAddr().String())

	if !ok {
		s.SendMessage(codes.SENDTYPE_CLIENT, s.Connlist[convert.ToString(UserId)], TypeMessage{"message", "用户id不存在，无法绑定"}, c, codes.MESSAGETYPE_TEXT)
	}
	//用户数据
	UserInfo := s.UserService.Repository.GetUserByID(UserId.(uint))
	user := &User{
		Conn:     conn,
		Name:     UserInfo.Phone,
		Id:       convert.ToString(UserId),
		Avator:   "string",
		To_id:    "string",
		group_id: "string",
		Failtime: 0,
	}
	s.register <- user

	for {
		//接受消息
		var receive []byte
		var recevString string
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		recevString = string(receive)
		log.Println("客户端:", recevString)

		var typeMsg TypeMessage
		json.Unmarshal(receive, &typeMsg)
		if typeMsg.Type == nil || typeMsg.Data == nil {
			log.Println("消息格式错误：", recevString)
			s.SendMessage(codes.SENDTYPE_CLIENT, user, TypeMessage{"message", "消息格式错误"}, c, codes.MESSAGETYPE_TEXT)
			continue
		}
		s.SendMessage(codes.SENDTYPE_SERVER, user, typeMsg, c, messageType)
	}
}
func (s *ServerManager) OnMessage() {
	for {
		select {
		case message := <-s.ClientMessage:
			log.Println("发送客户端:", string(message.content))
			s.Lock.Lock()
			err := message.user.Conn.WriteMessage(websocket.TextMessage, message.content)
			if err != nil {
				message.user.Failtime++
				if message.user.Failtime >= 3 {
					//如果发送失败次数超过指定次数则将该用户移除在线列表
					delete(s.Connlist, message.user.Id)
				} else {
					s.Connlist[message.user.Id] = message.user
				}
			} else {
				message.user.Failtime = 0
				s.Connlist[message.user.Id] = message.user
			}
			s.Lock.Unlock()
		case message := <-s.ServerMessage:
			var typeMsg TypeMessage
			json.Unmarshal(message.content, &typeMsg)
			conn := message.user.Conn
			msgType := typeMsg.Type.(string)
			log.Println("服务端接收:", string(message.content))

			switch msgType {
			//心跳
			case "ping":
				msg := TypeMessage{
					Type: "pong",
				}
				str, _ := json.Marshal(msg)
				s.Lock.Lock()
				conn.WriteMessage(websocket.TextMessage, str)
				s.Lock.Unlock()
			}
		case user := <-s.register:
			log.Println("注册用户:", user)
			s.Len++
			s.Connlist[user.Id] = user
			//注册成功向客户端发送个初始化信息

		}

	}

}
func (s *ServerManager) TestChatclient() {
	for {
		time.Sleep(time.Second * 25)
		for key, otherKefu := range s.Connlist {
			str := "发送给：" + key + "；" + time.Now().Format("2006-01-02 15:04:05")
			s.SendMessage(codes.SENDTYPE_CLIENT, otherKefu, TypeMessage{"message", str}, &gin.Context{}, codes.MESSAGETYPE_TEXT)
		}
	}
}
func (s *ServerManager) SendMessage(sendtype int, user *User, content TypeMessage, c *gin.Context, messageType int) {
	contentstr, _ := json.Marshal(content)
	if sendtype == codes.SENDTYPE_SERVER {
		s.ServerMessage <- &Message{
			user:        user,
			content:     contentstr,
			context:     c,
			messageType: messageType,
		}
	} else {
		s.ClientMessage <- &Message{
			user:        user,
			content:     contentstr,
			context:     c,
			messageType: messageType,
		}
	}

}
