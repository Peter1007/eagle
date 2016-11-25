package model

import (
	"eagle/core"
	"fmt"
)

type MemberModel struct {
}

func (this *MemberModel) GetList(reqNo string, statusId int64) string {
	data := map[string]interface{}{
		"status_id": statusId,
	}

	memberListRetStr := mq.CallService(reqNo, "user", "member/get_list", data)

	params := []core.CallParam{
		core.CallParam{Module: "user", Route: "member/get_list", Data: data},
		core.CallParam{Module: "base", Route: "member_type/get_list", Data: data},
	}
	ret, err := mq.CallMultiService(reqNo, params)
	fmt.Printf("%#v %#v\n", ret, err)

	return memberListRetStr
}
