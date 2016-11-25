package model

import (
	"eagle/core"
	"eagle/statistic/data"
	"time"
)

const (
	TABLE_PRE = "reqlog_"
)

type ReqLogModel struct {
}

var reqLogData *data.ReqLogData

func (this *ReqLogModel) Add(reqMeta core.ReqMeta, logData map[string]interface{}) error {
	reqTime := reqMeta.ReqTime
	nowTime := time.Now().Unix()
	if reqTime < nowTime-86400 {
		reqTime = nowTime
	}

	logTable := TABLE_PRE + time.Unix(reqTime, 0).Format("20060102")

	data := []interface{}{reqMeta.ReqNo, reqMeta.Step, reqMeta.Source, reqMeta.ReqTime, logData["module"], logData["route"],
		reqMeta.RemoteIp, logData["req_data"], logData["req_size"], logData["cost_time"], logData["code"], logData["msg"],
		logData["resp_data"], logData["resp_size"]}

	return reqLogData.Add(logTable, data...)
}
