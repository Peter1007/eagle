package model

type MemberTypeModel struct {
}

func (this *MemberTypeModel) GetList(reqNo string, statusId int64) (map[string]string, error) {
	data := map[string]interface{}{
		"status_id": statusId,
	}

	memberTypeList, err := mq.CallService(reqNo, "base", "member_type/get_list", data)
	if err != nil {
		return nil, err
	} else {
		return memberTypeList.(map[string]string), nil
	}
}
