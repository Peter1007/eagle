package core

import (
	"os"
	"path/filepath"
)

var appPath string

func init() {

	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	logger = GetLogInstance()

	routeCore = GetRouteCoreInstance()

	mq = GetMqInstance()
	mq.initConfig()
	mq.initConnsAndChannels()
}

func GetAppPath() string {
	return appPath
}
