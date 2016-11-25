package hander

import (
	"eagle/business/base/model"
)

type MemberTypeHander struct {
}

var memberTypeModel *model.MemberTypeModel

func (self *MemberTypeHander) GetList(reqNo string, statusId int64) (map[string]string, error) {
	return memberTypeModel.GetList(reqNo, statusId)
}
