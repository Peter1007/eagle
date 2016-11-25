package model

import (
	"eagle/service/user/data"
)

type MemberModel struct {
}

var memberData *data.MemberData

func (this *MemberModel) GetList(statusId int64) ([]map[string]string, error) {
	ret, err := memberData.GetList(statusId)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
