package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WebsocketClient struct {
	scheme      string
	siteName    string
	port        int
	path        string
	conn        *websocket.Conn
	close       chan struct{}
	closeOnce   sync.Once
	onCloseFunc func()
}

func NewWebsocketClient(scheme, siteName string, port int, path string) *WebsocketClient {
	client := &WebsocketClient{
		scheme:   scheme,
		siteName: siteName,
		port:     port,
		path:     path,
		close:    make(chan struct{}),
	}
	return client
}

func (client *WebsocketClient) Connect(header http.Header) error {
	u := url.URL{Scheme: client.scheme, Host: fmt.Sprintf("%s:%d", client.siteName, client.port), Path: client.path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	client.conn = c
	return nil
}

func (client *WebsocketClient) SetTickerFunc(f func() bool, d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-client.close:
				return
			case <-ticker.C:
				if !f() {
					return
				}
			}
		}
	}()
}

func (client *WebsocketClient) SendMsg(msg []byte) error {
	err := client.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

func (client *WebsocketClient) OnReceive(f func([]byte) bool) {
	go func() {
		defer client.Close()
		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				return
			}
			if !f(msg) {
				return
			}
		}
	}()
}

func (client *WebsocketClient) Close() {
	client.closeOnce.Do(func() {
		close(client.close)
		if client.conn != nil {
			_ = client.conn.Close()
		}
		if client.onCloseFunc != nil {
			client.onCloseFunc()
		}
	})
}

func (client *WebsocketClient) OnClose(f func()) {
	client.onCloseFunc = f
}
