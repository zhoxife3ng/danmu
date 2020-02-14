package message

import (
	"regexp"
)

const (
	RecvTypeGift         = "spbc"        // 全站礼物广播
	RecvTypeChatMsg      = "chatmsg"     // 弹幕消息
	RecvTypeUserEnter    = "uenter"      // 进入房间
	RecvTypeShareRoom    = "srres"       // 分享房间
	RecvTypeUserLevelUp  = "upgrade"     // 用户等级
	RecvTypeSuperChatMsg = "ssd"         // 超级弹幕
	RecvTypeBanned       = "newblackres" // 禁言
	RecvTypeError        = "error"       // 自定义error
)

var (
	RegexpUid, _              = regexp.Compile("/uid@=([0-9]+)/")
	RegexpType, _             = regexp.Compile("type@=(.*?)/")
	RegexpTypeGift, _         = regexp.Compile("/sn@=(.*?)/dn@=(.*?)/gn@=(.*?)/gc@=(.*?)/")
	RegexpTypeChatMsg, _      = regexp.Compile("/nn@=([^/]*?)/txt@=([^/]*?)/")
	RegexpTypeUserEnter, _    = regexp.Compile("/nn@=(.*?)/")
	RegexpTypeShareRoom, _    = regexp.Compile("/nickname@=(.*?)/")
	RegexpTypeUserLevelUp, _  = regexp.Compile("/nn@=(.*?)/level@=(.*?)/")
	RegexpTypeSuperChatMsg, _ = regexp.Compile("/content@=(.*?)/")
	RegexpTypeBanned, _       = regexp.Compile("/snic@=(.*?)/dnic@=(.*?)/")
)

type RecvMsg struct {
	Type string
	Data []string
	Uid  int64
}
