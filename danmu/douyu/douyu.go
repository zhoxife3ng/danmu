package douyu

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/x554462/danmu/client"
	. "github.com/x554462/danmu/danmu/douyu/message"
	"github.com/x554462/weuse/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DanmuDouyu struct {
	wsClient    *client.WebsocketClient
	recvChannel chan *RecvMsg
	closeChan   chan struct{}
}

func (dm *DanmuDouyu) Connect() (*client.WebsocketClient, error) {
	log.Println("连接斗鱼弹幕服务器")
	if err := checkRoom(); err != nil {
		return nil, err
	}
	wsClient := client.NewWebsocketClient(Scheme, SiteName, Port, "")
	header := http.Header{}
	header.Set("Origin", "https://www.douyu.com")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")
	if err := wsClient.Connect(header); err != nil {
		log.Println("dial:", err)
		return nil, err
	}
	dm.wsClient = wsClient
	return wsClient, nil
}

func (dm *DanmuDouyu) Login() error {
	log.Println("登录斗鱼房间")
	if err := dm.wsClient.SendMsg(NewSendMsg(SendTypeLoginRoom, DefaultRoomId).PackMsg()); err != nil {
		return err
	}
	if err := dm.wsClient.SendMsg(NewSendMsg(SendTypeJoinRoom, DefaultRoomId).PackMsg()); err != nil {
		return err
	}
	return nil
}

func (dm *DanmuDouyu) GetTickerTime() time.Duration {
	return 45 * time.Second
}

func (dm *DanmuDouyu) TickerFunc() bool {
	if err := dm.wsClient.SendMsg(NewSendMsg(SendTypeKeepLive, "").PackMsg()); err != nil {
		dm.wsClient.Close()
		return false
	}
	return true
}

func (dm *DanmuDouyu) OnReceive(b []byte) bool {
	msg := string(b)
	match := RegexpType.FindStringSubmatch(msg)
	if len(match) < 1 {
		return true
	}
	typo := match[1]

	message := &RecvMsg{Type: typo}

	match = RegexpUid.FindStringSubmatch(msg)
	if len(match) > 1 {
		if uid, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			message.Uid = uid
		}
	}
	switch typo {
	case RecvTypeGift:
		if m := RegexpTypeGift.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeChatMsg:
		if m := RegexpTypeChatMsg.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeUserEnter:
		if m := RegexpTypeUserEnter.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeShareRoom:
		if m := RegexpTypeShareRoom.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeUserLevelUp:
		if m := RegexpTypeUserLevelUp.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeSuperChatMsg:
		if m := RegexpTypeSuperChatMsg.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeBanned:
		if m := RegexpTypeBanned.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	default:
		return true
	}
	dm.recvChannel <- message
	return true
}

func (dm *DanmuDouyu) SetMsgChannel(msgChannel chan string) {
	dm.closeChan = make(chan struct{}, 0)
	dm.recvChannel = make(chan *RecvMsg, 10)
	go func() {
		for {
			select {
			case data := <-dm.recvChannel:
				if data.Type == RecvTypeChatMsg && len(data.Data) > 1 {
					msgChannel <- data.Data[1]
				}
			case <-dm.closeChan:
				return
			}
		}
	}()

}

func (dm *DanmuDouyu) OnClose() {
	log.Println("断开斗鱼连接")
	close(dm.closeChan)
}

func (dm *DanmuDouyu) Pause() {
	<-dm.closeChan
}

func checkRoom() error {
	log.Println("检查斗鱼房间")
	resBody, err := client.HttpReq(client.HttpMethodGet, fmt.Sprintf(RoomInfoUrl, DefaultRoomId), nil)
	if err != nil {
		log.Println("http:", err)
		return err
	}
	var resMap map[string]interface{}
	err = utils.JsonDecodeWithByte(resBody, &resMap)
	if err != nil {
		log.Println("json decode:", err)
		return err
	}
	errMsg := ""
	if resErr, ok := (resMap["error"]).(float64); ok {
		if resErr != 0 {
			errMsg = "房间号不存在"
		} else {
			if data, ok := resMap["data"].(map[string]interface{}); ok {
				if roomStatus, ok := data["room_status"].(string); ok {
					if roomStatus == "2" {
						errMsg = "主播未开播"
					}
				} else {
					errMsg = "room status error"
				}
			} else {
				errMsg = "data error"
			}
		}
	} else {
		errMsg = "json error"
	}
	if errMsg != "" {
		log.Println("http:", errMsg)
		return errors.New(errMsg)
	}
	return nil
}
