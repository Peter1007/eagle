package main

import (
	"eagle/business/user/route"
	"eagle/core"
)

const (
	ENV_FILE = "env.yml"
)

type EnvConfig struct {
	Module string `json:"module"`
}

func main() {
	mq := core.GetMqInstance()
	defer mq.Close()

	logger := core.GetLogInstance()
	defer logger.Close()

	envConfig := EnvConfig{}
	err := core.SetConfig(ENV_FILE, &envConfig)
	if err != nil {
		logger.LogPanic(err)
	}

	mq.SetModule(envConfig.Module)

	route.InitRoute()

	mq.Run()
}
