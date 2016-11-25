package data

import (
	"eagle/core/db"
	"strconv"
)

type MemberTypeData struct {
}

func (this *MemberTypeData) GetList(statusId int64) ([]map[string]string, error) {
	sql := "SELECT `type_id`, `type_name`, `status` FROM `member_type`"
	if statusId > 0 {
		sql = sql + " WHERE `status`='" + strconv.FormatInt(statusId, 10) + "'"
	}

	dbConn, err := db.GetConnection(DB_BASE)
	if err != nil {
		return nil, err
	}

	return dbConn.Select(sql)
}
