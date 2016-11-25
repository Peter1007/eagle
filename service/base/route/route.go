package route

import (
	"eagle/core"
	"eagle/service/base/hander"
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
	cityHander := &hander.CityHander{}
	cityReflect := reflect.ValueOf(&cityHander).Elem()
	routeCore.AddRouteHander("city/get_list", cityReflect.MethodByName("GetList"))

	memberTypeHander := &hander.MemberTypeHander{}
	memberTypeReflect := reflect.ValueOf(&memberTypeHander).Elem()
	routeCore.AddRouteHander("member_type/get_list", memberTypeReflect.MethodByName("GetList"))
}
