package data

import (
	"eagle/core/db"
	"strconv"
)

type CityData struct {
}

func (this *CityData) GetList(statusId int64) ([]map[string]string, error) {
	sql := "SELECT `city_id`, `city_name` FROM `city`"
	if statusId > 0 {
		sql = sql + " WHERE `status`='" + strconv.FormatInt(statusId, 10) + "'"
	}

	dbConn, err := db.GetConnection(DB_BASE)
	if err != nil {
		return nil, err
	}

	return dbConn.Select(sql)
}
