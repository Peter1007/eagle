package hander

import (
	"eagle/service/base/model"
)

type MemberTypeHander struct {
}

var memberTypeModel *model.MemberTypeModel

func (self *MemberTypeHander) GetList(statusId int64) (map[string]string, error) {
	return memberTypeModel.GetList(statusId)
}
