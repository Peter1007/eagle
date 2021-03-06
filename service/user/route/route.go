package route

import (
	"eagle/core"
	"eagle/service/user/hander"
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
	memberHander := &hander.MemberHander{}
	memberReflect := reflect.ValueOf(&memberHander).Elem()

	routeCore.AddRouteHander("member/get_list", memberReflect.MethodByName("GetList"))
}
