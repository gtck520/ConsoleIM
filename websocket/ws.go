package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/gtck520/ConsoleIM/common/http"
	"github.com/gtck520/ConsoleIM/common/logger"
)

type Ws struct {
	Conn *websocket.Conn
	Log  logger.Logger
}

func (w *Ws) Connect(token string) {
	zap := logger.Logger{}
	zap.Init()
	w.Log = zap
	c, _, err := websocket.DefaultDialer.Dial(http.Websocket_Url+"?X-Token="+token, nil)
	if err != nil {
		zap.Error("dial:", err)

	}
	w.Conn = c
}
