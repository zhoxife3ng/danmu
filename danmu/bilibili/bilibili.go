package bilibili

import (
	"encoding/binary"
	"github.com/x554462/danmu/client"
	. "github.com/x554462/danmu/danmu/bilibili/message"
	"github.com/x554462/danmu/danmu/bilibili/util"
	"github.com/x554462/weuse/utils"
	"log"
	"net/http"
	"time"
)

type DanmuBilibili struct {
	wsClient    *client.WebsocketClient
	config      *Config
	recvChannel chan string
	closeChan   chan struct{}
}

func (dm *DanmuBilibili) Connect() (*client.WebsocketClient, error) {
	log.Println("连接B站弹幕服务器")
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	wsClient := client.NewWebsocketClient(config.Scheme, config.SiteName, config.Port, config.Path)

	header := http.Header{}
	header.Set("Origin", "https://live.bilibili.com")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")
	if err := wsClient.Connect(header); err != nil {
		log.Fatalln("dial:", err)
		return nil, err
	}
	dm.config = config
	dm.wsClient = wsClient
	return wsClient, nil
}

func (dm *DanmuBilibili) Login() error {
	log.Println("登录B站房间")
	msg, err := NewMsgLogin(0, DefaultRoomId, ProtocolVersionDeflate, 2, "web", "1.10.1", dm.config.Key).PackMsg(ProtocolVersionNormal)
	if err != nil {
		return err
	}
	return dm.wsClient.SendMsg(msg)
}

func (dm *DanmuBilibili) GetTickerTime() time.Duration {
	return 30 * time.Second
}

func (dm *DanmuBilibili) TickerFunc() bool {
	msg, err := NewMsgKeepLive("[object Object]").PackMsg(ProtocolVersionNormal)
	if err == nil {
		if err = dm.wsClient.SendMsg(msg); err == nil {
			return true
		}
	}
	dm.wsClient.Close()
	return false
}

func (dm *DanmuBilibili) OnReceive(b []byte) bool {

	packetLen := binary.BigEndian.Uint32(b[HeadMsgLenOffset : HeadMsgLenOffset+4])
	headLen := binary.BigEndian.Uint16(b[HeadHeadLenOffset : HeadHeadLenOffset+2])
	protocolVer := binary.BigEndian.Uint16(b[HeadProtocolVerOffset : HeadProtocolVerOffset+2])
	op := binary.BigEndian.Uint32(b[HeadOpOffset : HeadOpOffset+4])
	//seq := binary.BigEndian.Uint32(b[headSeqOffset : headSeqOffset+4])
	if op == WsOpMessage {
		msgLen := len(b)
		for offset := 0; offset < msgLen; offset += int(packetLen) {
			data := ""
			packetLen = binary.BigEndian.Uint32(b[offset : offset+4])
			headLen = binary.BigEndian.Uint16(b[offset+HeadHeadLenOffset : offset+HeadHeadLenOffset+2])
			protocolVer = binary.BigEndian.Uint16(b[offset+HeadProtocolVerOffset : offset+HeadProtocolVerOffset+2])
			msg := b[offset+int(headLen) : offset+int(packetLen)]
			switch protocolVer {
			case ProtocolVersionNormal:
				if msgStr := RegexpJson.FindString(string(msg)); msgStr != "" {
					data = msgStr
				}
			case ProtocolVersionDeflate:
				msg, err := util.GzInflate(msg[2:])
				if err == nil {
					if msgStr := RegexpJson.FindString(string(msg)); msgStr != "" {
						data = msgStr
					}
				}
			default:
				data = ""
			}
			if data != "" {
				dm.recvChannel <- data
			}
		}
	}
	return true
}

func (dm *DanmuBilibili) SetMsgChannel(msgChannel chan string) {
	dm.closeChan = make(chan struct{}, 0)
	dm.recvChannel = make(chan string, 10)
	var danmuMsg struct {
		Cmd  string      `json:"cmd"`
		Info interface{} `json:"info"`
	}
	go func() {
		for {
			select {
			case data := <-dm.recvChannel:
				if err := utils.JsonDecode(data, &danmuMsg); err == nil && danmuMsg.Cmd == "DANMU_MSG" {
					if info, ok := danmuMsg.Info.([]interface{}); ok {
						if len(info) > 1 {
							if msg, ok := info[1].(string); ok {
								msgChannel <- msg
							}
						}
					}
				}
			case <-dm.closeChan:
				return
			}
		}
	}()
}

func (dm *DanmuBilibili) OnClose() {
	log.Println("断开B站连接")
	close(dm.closeChan)
}

func (dm *DanmuBilibili) Pause() {
	<-dm.closeChan
}
