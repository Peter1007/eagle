package core

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"math/rand"
)

const (
	ERRNO_ERROR  = 10002
	CODE_SUCCESS = 0
	MSG_SUCCESS  = "success"
	CONFIG_DIR   = "/config/"
)

type ResponseBody struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SetConfig(configFile string, out interface{}) error {
	filename := appPath + CONFIG_DIR + configFile

	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Read mq config file error. %s", err.Error()))
	}

	yaml.Unmarshal(buffer, out)

	return nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)

	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

func ResponseOk(data interface{}) ResponseBody {
	return ResponseBody{
		Code: CODE_SUCCESS,
		Msg:  MSG_SUCCESS,
		Data: data,
	}
}

func ResponseFail(code int, err error) ResponseBody {
	if code == 0 {
		code = ERRNO_ERROR
	}

	return ResponseBody{
		Code: code,
		Msg:  err.Error(),
		Data: map[string]string{},
	}
}
