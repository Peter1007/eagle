package hander

import (
	"eagle/service/base/model"
)

type CityHander struct {
}

var cityModel *model.CityModel

func (self *CityHander) GetList(statusId int64) (map[string]string, error) {
	return cityModel.GetList(statusId)
}
