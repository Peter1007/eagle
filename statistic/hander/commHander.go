package hander

import (
	"eagle/core"
)

var logger *core.Logger

func init() {
	logger = core.GetLogInstance()
}
