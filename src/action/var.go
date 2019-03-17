package action

import (
	"model"

	jsoniter "github.com/json-iterator/go"
)

// JSON JSON
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

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
