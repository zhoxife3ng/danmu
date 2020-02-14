package danmu

import (
	"fmt"
	"github.com/x554462/danmu/client"
	"github.com/x554462/danmu/danmu/bilibili"
	"github.com/x554462/danmu/danmu/douyu"
	"time"
)

type Danmu interface {
	Connect() (*client.WebsocketClient, error)
	Login() error
	GetTickerTime() time.Duration
	TickerFunc() bool
	OnReceive(b []byte) bool
	OnClose()
	SetMsgChannel(msgChannel chan string)
	Pause()
}

var (
	msgChannel = make(chan string, 10)
)

var DanmuPro = []Danmu{
	&douyu.DanmuDouyu{},
	&bilibili.DanmuBilibili{},
}

func Run() {
	go handle()
	for _, danmu := range DanmuPro {
		go runDanmu(danmu)
	}
}

func handle() {
	for {
		select {
		case msg := <-msgChannel:
			fmt.Println(msg)
		}
	}
}

const loopDuration = time.Second << 5

func runDanmu(danmu Danmu) {
	var (
		wsClient *client.WebsocketClient
		err      error
		timer    = time.NewTimer(0)
	)
	defer timer.Stop()
	for {
		resetTime := loopDuration
		select {
		case <-timer.C:
			wsClient, err = danmu.Connect()
			if err == nil {
				if err = danmu.Login(); err == nil {
					tickerTime := danmu.GetTickerTime()
					if tickerTime > 0 {
						wsClient.SetTickerFunc(danmu.TickerFunc, tickerTime)
					}
					danmu.SetMsgChannel(msgChannel)
					wsClient.OnReceive(danmu.OnReceive)
					wsClient.OnClose(danmu.OnClose)
					danmu.Pause()
					resetTime = 3 * time.Second
				}
			}
		}
		timer.Reset(resetTime)
	}
}
