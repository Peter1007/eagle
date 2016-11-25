package model

import (
	"eagle/service/base/data"
)

type MemberTypeModel struct {
}

var memberTypeData *data.MemberTypeData

func (this *MemberTypeModel) GetList(statusId int64) (map[string]string, error) {
	result, err := memberTypeData.GetList(statusId)
	if err != nil {
		return nil, err
	}

	ret := map[string]string{}
	for _, value := range result {
		ret[value["type_id"]] = value["type_name"]
	}

	return ret, nil
}
