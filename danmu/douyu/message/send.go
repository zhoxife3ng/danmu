package message

import (
	"encoding/binary"
	"strings"
)

const (
	SendTypeLoginRoom = "type@=loginreq/roomid@=:msg/"
	SendTypeKeepLive  = "type@=mrkl/"
	SendTypeJoinRoom  = "type@=joingroup/rid@=:msg/gid@=-9999/"
	SendTypeLogout    = "type@=logout/"
)

type SendMsg struct {
	Type  string
	Param string
}

func NewSendMsg(typo, param string) *SendMsg {
	return &SendMsg{Type: typo, Param: param}
}

func (sm *SendMsg) PackMsg() []byte {
	msg := strings.Replace(sm.Type, ":msg", sm.Param, -1)

	length := 10 + len(msg)                // 4(code)+4(magic)+x(msg)+2(end)
	var result = make([]byte, 4, 4+length) // 4(length) + length
	binary.LittleEndian.PutUint32(result, uint32(length))

	result = append(result, result...)
	result = append(result, 0xb1, 0x02, 0x00, 0x00)
	result = append(result, []byte(msg)...)
	return append(result, 0x00, 0x00)
}
