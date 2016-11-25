package model

import (
	"eagle/core"
)

var mq *core.Mq

func init() {
	mq = core.GetMqInstance()
}
