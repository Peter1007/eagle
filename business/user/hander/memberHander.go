package hander

import (
	"eagle/business/user/model"
)

type MemberHander struct {
}

var memberModel *model.MemberModel

func (self *MemberHander) GetList(reqNo string, statusId int64) string {
	return memberModel.GetList(reqNo, statusId)
}
