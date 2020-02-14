package message

import (
	"regexp"
)

const (
	HeadMsgLenOffset      = 0
	HeadHeadLenOffset     = 4
	HeadProtocolVerOffset = 6
	HeadOpOffset          = 8
	HeadSeqOffset         = 12
)

var (
	RegexpJson, _ = regexp.Compile("{.*}")
)
