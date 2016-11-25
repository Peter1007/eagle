package data

import (
	"eagle/core/db"
)

type ReqLogData struct {
}

func (this *ReqLogData) Add(table string, data ...interface{}) error {
	sql := `INSERT INTO ` + table + ` (req_no, step, source, req_time, module, route, remote_ip, req_data, req_data_size,
	cost_time, resp_code, resp_msg, resp_data, resp_data_size) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	dbConn, err := db.GetConnection(DB_REQ_LOG)
	if err != nil {
		return err
	}

	_, err = dbConn.Insert(sql, data...)

	return err
}
