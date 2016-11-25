package model

type CityModel struct {
}

func (this *CityModel) GetList(reqNo string, statusId int64) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"status_id": statusId,
	}

	cityList, err := mq.CallService(reqNo, "base", "city/get_list", data)
	if err != nil {
		return nil, err
	} else {
		return cityList.(map[string]interface{}), nil
	}
}
