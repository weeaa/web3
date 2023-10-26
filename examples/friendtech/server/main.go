package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/friendtech/watcher"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Shows how to create a websocket server & broadcast Friend Tech data.
func main() {
	server := SocketServer{
		Clients:    make(map[*websocket.Conn]*socketClient),
		httpClient: &http.Client{},
		eventChan:  make(chan EventChannelItem),
	}

	go server.StartWsServer()
	pgConn, err := db.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	discBot, err := bot.New(pgConn)
	if err != nil {
		log.Fatal(err)
	}

	if err = discBot.Start(); err != nil {
		log.Fatal(err)
	}

	ftWatcher, err := watcher.NewFriendTech(pgConn, discBot, "proxies.txt", os.Getenv("NODE_WSS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	data := <-ftWatcher.OutStreamData
	_ = data
	server.eventChan <- EventChannelItem{}
}

const listenPath = "/var/run/friendtech.sock"

type socketClient struct {
	Users    []string
	messages chan map[string]interface{}

	close bool
}

type socketRequest struct {
	Event string      `json:"event"`
	Data  interface{} `json:"trade,omitempty"`
}

type SocketServer struct {
	Clients    map[*websocket.Conn]*socketClient
	upgrader   websocket.Upgrader
	httpClient *http.Client

	eventChan chan EventChannelItem
}

type EventChannelItem struct {
	Address string
}

func (s *SocketServer) DeleteClient(c *websocket.Conn) {
	if s.Clients[c].close {
		return
	}
	if err := c.Close(); err != nil {
		return
	}
	s.Clients[c].close = true
	delete(s.Clients, c)
}

func (s *SocketServer) ReadClient(c *websocket.Conn) {
	var err error
	for {
		if s.Clients[c].close {
			return
		}

		var msg map[string]interface{}
		err = c.ReadJSON(&msg)
		if err != nil {
			return
		}
		s.Clients[c].messages <- msg
	}
}

// Advised for authenticating to add a signature header like 'x-sign'
// to make it safer and authenticate users.
func (s *SocketServer) HandleFunc(wh http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		wh.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users := r.Header.Get("x-users")
	if users == "[]" || users == "" {
		wh.WriteHeader(http.StatusBadRequest)
		wh.Write([]byte(""))
		return
	}

	var clientUsers []string
	if err := json.Unmarshal([]byte(users), &clientUsers); err != nil {
		wh.WriteHeader(http.StatusBadRequest)
		return
	}

	clientConn, err := s.upgrader.Upgrade(wh, r, nil)
	if err != nil {
		return
	}

	s.Clients[clientConn] = &socketClient{Users: clientUsers, messages: make(chan map[string]interface{}, 300)}

	go s.ReadClient(clientConn)
}

func (s *SocketServer) StartWsServer() {
	s.Clients = map[*websocket.Conn]*socketClient{}
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}

	s.httpClient = &http.Client{Transport: &http.Transport{}}

	http.HandleFunc("/", s.HandleFunc)

	go func() {
		_ = os.Remove(listenPath)
		conn, err := net.Listen("unix", listenPath)
		if err != nil {
			panic(err)
		}
		_ = os.Chmod(listenPath, 0777)
		err = http.Serve(conn, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go s.startBroadcaster()
	go s.pingSender()
}

// pingSender sends frequently 'ping' messages as a heartbeat message.
func (s *SocketServer) pingSender() {
	pingMsg := socketRequest{
		Event: "ping",
	}

	for {
		time.Sleep(time.Minute)
		for client := range s.Clients {
			go func(ws *websocket.Conn) {
				if err := ws.WriteJSON(pingMsg); err != nil {
					s.DeleteClient(ws)
					return
				}

				for len(s.Clients[ws].messages) > 0 {
					<-s.Clients[ws].messages
				}

				select {
				case <-time.NewTimer(time.Second * 5).C:
					s.DeleteClient(ws)
				case <-s.Clients[ws].messages:
				}

			}(client)
		}
	}
}

// startBroadcaster broadcasts data to users subscribed.
// as soon as s.eventChan receives data, it is sent to users.
// todo update with accurate trade data
func (s *SocketServer) startBroadcaster() {
	for {
		event := <-s.eventChan
		for socketConn, client := range s.Clients {
			go func(ws *websocket.Conn, client *socketClient, userItem EventChannelItem) {
				for _, user := range client.Users {
					if user == userItem.Address {
						if err := ws.WriteJSON(socketRequest{Event: "", Data: "tradeItem.Trade"}); err != nil {
							s.DeleteClient(ws)
							return
						}
					}
				}
			}(socketConn, client, event)
		}
	}
}
