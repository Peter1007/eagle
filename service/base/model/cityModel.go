package model

import (
	"eagle/service/base/data"
)

type CityModel struct {
}

var cityData *data.CityData

func (this *CityModel) GetList(statusId int64) (map[string]string, error) {
	result, err := cityData.GetList(statusId)
	if err != nil {
		return nil, err
	}

	ret := map[string]string{}
	for _, value := range result {
		ret[value["city_id"]] = value["city_name"]
	}

	return ret, nil
}
