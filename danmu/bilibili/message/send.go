package message

import (
	"encoding/binary"
	"github.com/x554462/danmu/danmu/bilibili/util"
	"github.com/x554462/weuse/utils"
)

const (
	ProtocolVersionNormal  uint16 = 0
	ProtocolVersionDeflate uint16 = 2
)

const (
	WsOpHeartbeat            = 2
	WsOpHeartbeatReply       = 3
	WsOpMessage              = 5
	WsOpUserAuthentication   = 7
	WsOpConnectSuccess       = 8
	WsHeaderDefaultVersion   = 1
	WsHeaderDefaultOperation = 1
	WsHeaderDefaultSequence  = 1
)

type SendMsg interface {
	PackMsg(protocol uint8) ([]byte, error)
}

type msgLogin struct {
	Uid       int    `json:"uid"`
	RoomId    int    `json:"roomid"`
	ProtoVer  uint16 `json:"protover"`
	Platform  string `json:"platform"`
	ClientVer string `json:"clientver"`
	Type      int    `json:"type"`
	Key       string `json:"key"`
}

func NewMsgLogin(uid, roomId int, protoVer uint16, typo int, platform, clientVer, key string) *msgLogin {
	return &msgLogin{
		Uid:       uid,
		RoomId:    roomId,
		ProtoVer:  protoVer,
		Platform:  platform,
		ClientVer: clientVer,
		Type:      typo,
		Key:       key,
	}
}

func (ml *msgLogin) PackMsg(protocol uint16) (b []byte, err error) {
	var bodyByte []byte
	if protocol == ProtocolVersionNormal {
		bodyByte, err = utils.JsonEncodeByte(ml)
	} else if protocol == ProtocolVersionDeflate {
		b1, err := utils.JsonEncodeByte(ml)
		if err == nil {
			bodyByte, err = util.GzDeflate(b1)
		}
	}
	if err != nil {
		return
	}
	b = Pack(bodyByte, WsHeaderDefaultVersion, WsOpUserAuthentication, WsHeaderDefaultSequence)
	return
}

type msgKeepLive struct {
	Data string
}

func NewMsgKeepLive(data string) *msgKeepLive {
	return &msgKeepLive{Data: data}
}

func (mkl *msgKeepLive) PackMsg(protocol uint16) (b []byte, err error) {
	var bodyByte []byte
	if protocol == ProtocolVersionNormal {
		bodyByte = []byte(mkl.Data)
	} else if protocol == ProtocolVersionDeflate {
		bodyByte, err = util.GzDeflate([]byte(mkl.Data))
	}
	if err != nil {
		return
	}
	b = Pack(bodyByte, WsHeaderDefaultVersion, WsOpHeartbeat, WsHeaderDefaultSequence)
	return
}

func Pack(bodyByte []byte, version uint16, op, seq uint32) (b []byte) {
	headLen := 16
	bodyLen := len(bodyByte)
	b = make([]byte, headLen, bodyLen+headLen)
	b = append(b, bodyByte...)
	binary.BigEndian.PutUint32(b[0:], uint32(bodyLen+headLen))
	binary.BigEndian.PutUint16(b[4:], uint16(headLen))
	binary.BigEndian.PutUint16(b[6:], version)
	binary.BigEndian.PutUint32(b[8:], op)
	binary.BigEndian.PutUint32(b[12:], seq)
	return
}
