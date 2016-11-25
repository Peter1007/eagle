package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	CONFIG_PARAM_PATH = "param/%s.yml"
	INTEGER           = "integer"
	STRING            = "string"
)

type ParamConfig struct {
	Required     int
	DefaultValue interface{}
	ValueType    string
	Min          int
	Max          int
}

type Param struct{}

var param *Param

func GetParamInstance() *Param {
	if param == nil {
		param = &Param{}
	}

	return param
}

func (this *Param) GenReflectParam(route string, reqData map[string]interface{}) ([]reflect.Value, error) {
	reflectParams := []reflect.Value{}

	paramsConfig, err := this.getParamsConfig(route)
	if err != nil {
		return reflectParams, err
	}

	for key, config := range paramsConfig {
		if !this.checkRequired(&config, key, reqData) {
			return reflectParams, errors.New(key + " is required.")
		}

		paramValue := this.getParamValue(&config, key, reqData)

		if !this.checkType(&config, paramValue) {
			return reflectParams, errors.New(key + " is type error.")
		}

		if !this.checkRange(&config, paramValue) {
			return reflectParams, errors.New(key + " is too small or too large.")
		}

		if config.ValueType == INTEGER {
			value, _ := paramValue.(float64)
			reflectParams = append(reflectParams, reflect.ValueOf(int64(value)))
		} else if config.ValueType == STRING {
			value, _ := paramValue.(string)
			reflectParams = append(reflectParams, reflect.ValueOf(value))
		}

	}

	return reflectParams, err
}

func (this *Param) CheckMeta(reqMeta *ReqMeta, reqType int) error {
	if reqType == 0 {
		if len(reqMeta.ReqNo) < 16 {
			return errors.New("Meta's req_no error.")
		}

		if reqMeta.ReqTime < 1000000000 {
			return errors.New("Meta's req_time error.")
		}

		if len(reqMeta.Source) < 1 {
			return errors.New("Meta's source error.")
		}

		if reqMeta.Step < 0 {
			return errors.New("Meta's step error.")
		}
	}

	return nil
}

func (this *Param) getParamsConfig(route string) (map[string]ParamConfig, error) {
	routeSlice := strings.Split(route, "/")
	routeMain := routeSlice[0]
	routeSub := strings.Join(routeSlice[1:], "/")

	filename := fmt.Sprintf(CONFIG_PARAM_PATH, routeMain)

	configs := map[string]map[string]ParamConfig{}
	err := SetConfig(filename, &configs)
	if err != nil {
		return nil, err
	}

	return configs[routeSub], nil
}

func (this *Param) checkRequired(config *ParamConfig, key string, reqData map[string]interface{}) bool {
	if config.Required == 1 {
		_, ok := reqData[key]
		return ok
	} else {
		return true
	}
}

func (this *Param) getParamValue(config *ParamConfig, key string, reqData map[string]interface{}) interface{} {
	var paramValue interface{}

	value, ok := reqData[key]
	if ok {
		paramValue = value
	} else {
		paramValue = config.DefaultValue
	}

	return paramValue
}

func (this *Param) checkType(config *ParamConfig, paramValue interface{}) bool {
	ok := true

	if config.ValueType == INTEGER {
		_, ok = paramValue.(float64)
	} else if config.ValueType == STRING {
		_, ok = paramValue.(string)
	}

	return ok
}

func (this *Param) checkRange(config *ParamConfig, paramValue interface{}) bool {
	num := 0
	if config.ValueType == INTEGER {
		value, _ := paramValue.(float64)
		num = int(value)
	} else if config.ValueType == STRING {
		value, _ := paramValue.(string)
		num = len(value)
	}

	if num < config.Min || (0 < config.Max && config.Max < num) {
		return false
	}

	return true
}
