package core

import (
	"errors"
	"reflect"
)

type RouteCore struct{}

var routeCore *RouteCore
var routeMap map[string]reflect.Value

func GetRouteCoreInstance() *RouteCore {
	if routeCore == nil {
		routeCore = &RouteCore{}

		routeMap = map[string]reflect.Value{}
	}

	return routeCore
}

func (this *RouteCore) AddRouteHander(route string, hander reflect.Value) {
	routeMap[route] = hander
}

func (this *RouteCore) GetHander(route string) (reflect.Value, error) {
	hander, ok := routeMap[route]

	if ok {
		return hander, nil
	} else {
		return reflect.Value{}, errors.New("route: " + route + " is not exist.")
	}
}
