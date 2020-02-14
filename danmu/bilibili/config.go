package bilibili

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/x554462/danmu/client"
	"github.com/x554462/weuse/utils"
)

const (
	DefaultRoomId = 21859078
	ConfigInfoUrl = "https://api.live.bilibili.com/room/v1/Danmu/getConf?room_id=%d&platform=pc&player=web"
)

type Config struct {
	Scheme   string
	SiteName string
	Port     int
	Key      string
	Path     string
}

func NewConfig() (*Config, error) {
	body, err := client.HttpReq(client.HttpMethodGet, fmt.Sprintf(ConfigInfoUrl, DefaultRoomId), "")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		Message string `json:"message"`
		Data    struct {
			HosServerList []struct {
				Host    string `json:"host"`
				Port    int    `json:"port"`
				WssPort int    `json:"wss_port"`
				WsPort  int    `json:"ws_port"`
			} `json:"host_server_list"`
			Token string `json:"token"`
		} `json:"data"`
	}
	err = utils.JsonDecodeWithByte(body, &resp)
	if err != nil {
		return nil, err
	} else if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	} else if len(resp.Data.HosServerList) == 0 {
		return nil, errors.New("host list is nil")
	}
	return &Config{
		Scheme:   "wss",
		SiteName: resp.Data.HosServerList[0].Host,
		Port:     443,
		Key:      resp.Data.Token,
		Path:     "/sub",
	}, nil
}
