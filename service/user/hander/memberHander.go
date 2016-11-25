package hander

import (
	"eagle/core"
	"eagle/service/user/model"
)

type MemberHander struct {
}

var memberModel *model.MemberModel

func (self *MemberHander) GetList(statusId int64) string {
	data, err := memberModel.GetList(statusId)
	if err != nil {
		return core.ResponseFail(err)
	} else {
		return core.ResponseOk(data)
	}
}
