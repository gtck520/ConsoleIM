package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/gtck520/ConsoleIM/common/http"
	"github.com/gtck520/ConsoleIM/common/logger"
)

type Ws struct {
	Conn *websocket.Conn
}

func (w *Ws) Connect(token string) {
	zap := logger.Logger{}
	zap.Init()
	zap.Info("wsurl:", http.Websocket_Url+"?X-Token="+token)
	c, _, err := websocket.DefaultDialer.Dial(http.Websocket_Url+"?X-Token="+token, nil)
	if err != nil {
		zap.Error("dial:", err)

	}
	w.Conn = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				zap.Error("read:", err)
				return
			}
			zap.Infof("recv: %s", message)

		}
	}()
}
