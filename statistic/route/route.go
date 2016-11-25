package route

import (
	"eagle/core"
	"eagle/statistic/hander"
	"reflect"
)

var routeCore *core.RouteCore

func InitRoute() {
	routeCore = core.GetRouteCoreInstance()

	addRoute()
}

func GetHander(route string) (reflect.Value, error) {
	return routeCore.GetHander(route)
}

func addRoute() {
	reqLogHander := &hander.ReqLogHander{}
	reqLogReflect := reflect.ValueOf(&reqLogHander).Elem()

	routeCore.AddRouteHander("reqlog/add", reqLogReflect.MethodByName("Add"))
}
