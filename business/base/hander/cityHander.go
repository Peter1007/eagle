package hander

import (
	"eagle/business/base/model"
)

type CityHander struct {
}

var cityModel *model.CityModel

func (self *CityHander) GetList(reqNo string, statusId int64) (map[string]interface{}, error) {
	return cityModel.GetList(reqNo, statusId)
}
