package action

import (
	"model"

	jsoniter "github.com/json-iterator/go"
)

// JSON JSON
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

// routerNameSlice 去重用，router name 不能重复，否则只会留下第一个
var routerNameSlice []string

// routerPathSlice 去重用，router path 不能重复，否则只会留下第一个
var routerPathSlice []string

// RawRouterStruct RawRouterStruct
type RawRouterStruct struct {
	RouterName string `json:"name"`
	Method     string `json:"method"`
	RouterPath string `json:"path"`
	GroupName  string `json:"group"`
	URL        model.URLStruct
	Headers    []model.HeaderStruct
	Querys     []model.QueryStruct
	Files      []model.ModeDataFileStruct
	Texts      []model.ModeDataTextStruct
	Body       interface{}
}
