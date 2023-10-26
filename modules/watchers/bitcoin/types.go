package bitcoin

import (
	"github.com/gorilla/websocket"
)

const (
	wsBaseEndpoint = "wss://mempool.space/api/v1/ws"
)

type SocketClient struct {
	Conn *websocket.Conn
}
