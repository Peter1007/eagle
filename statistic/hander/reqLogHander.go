package hander

import (
	"eagle/core"
	"eagle/statistic/model"
)

type ReqLogHander struct {
}

var reqLogModel *model.ReqLogModel

func (self *ReqLogHander) Add(reqMeta core.ReqMeta, logData map[string]interface{}) {
	err := reqLogModel.Add(reqMeta, logData)
	if err != nil {
		logger.LogWarm(err.Error())
	}
}
