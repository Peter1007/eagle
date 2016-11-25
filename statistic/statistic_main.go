package main

import (
	"eagle/core"
	"eagle/core/db"
	"eagle/statistic/route"
)

func main() {
	mq := core.GetMqInstance()
	defer mq.Close()

	logger := core.GetLogInstance()
	defer logger.Close()

	route.InitRoute()

	err := db.InitDb()
	if err != nil {
		logger.LogPanic(err)
	}

	mq.Run()
}
