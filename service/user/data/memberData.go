package data

import (
	"eagle/core/db"
	"strconv"
)

type MemberData struct {
}

func (this *MemberData) GetList(statusId int64) ([]map[string]string, error) {
	sql := "SELECT `member_id`, `member_name`, `city_id`, `member_type_id`, `status`, `create_time` FROM `member`"
	if statusId > 0 {
		sql = sql + " WHERE `status`='" + strconv.FormatInt(statusId, 10) + "'"
	}

	dbConn, err := db.GetConnection(DB_USER)
	if err != nil {
		return nil, err
	}

	return dbConn.Select(sql)
}
