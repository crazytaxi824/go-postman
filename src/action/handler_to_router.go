package action

import (
	"errors"
	"log"
	"model"
	"strings"
)

// HandlersToRouters HandlersToRouters
func (router *RawRouterStruct) HandlersToRouters(handlers []string) {
	var handlerIndex []int
	for k, h := range handlers {
		if strings.Contains(h, "@pmHandler") {
			handlerIndex = append(handlerIndex, k)
		}
	}

	// handler name
	for i := 0; i < len(handlerIndex); i++ {
		tmp := strings.Split(handlers[handlerIndex[i]], "@pmHandler")
		if len(tmp) > 1 {
			handlerRef, err := ParsePMstructToJSONformat(strings.TrimSpace(tmp[1]))
			if err != nil {
				log.Println("warning: 格式错误 ——" + handlers[handlerIndex[i]])
				continue
			}

			data := make(map[string]string)
			err = JSON.UnmarshalFromString(handlerRef, &data)
			if err != nil {
				log.Println("warning: 格式错误 —— " + handlers[handlerIndex[i]])
				continue
			}

			if router.RouterName == data["name"] {
				if i > len(handlerIndex)-2 {
					for _, handler := range handlers[handlerIndex[i]+1:] {
						err = router.parseQueryBodyHeaders(handler)
						if err != nil {
							log.Println("warning: 格式错误 —— " + handler)
							continue
						}
					}
				} else {
					for _, handler := range handlers[handlerIndex[i]+1 : handlerIndex[i+1]] {
						err = router.parseQueryBodyHeaders(handler)
						if err != nil {
							log.Println("warning: 格式错误 —— " + handler)
							continue
						}
					}
				}
			}
		}
	}
	return
}

func (router *RawRouterStruct) parseQueryBodyHeaders(handler string) error {
	if strings.Contains(handler, "@pmQuery") {
		tmpQuery := strings.Split(handler, "@pmQuery")
		if len(tmpQuery) > 1 {

			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpQuery[1]))
			if err != nil {
				return err
			}

			dataQuery := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataQuery)
			if err != nil {
				return errors.New(handler + " —— 格式错误")
			}

			var query model.QueryStruct
			query.Key = dataQuery["key"]
			query.Value = dataQuery["value"]
			query.Description = dataQuery["desc"]
			router.Querys = append(router.Querys, query)
		}
	} else if strings.Contains(handler, "@pmBody") {
		tmpBody := strings.Split(handler, "@pmBody")
		if len(tmpBody) > 1 {
			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpBody[1]))
			if err != nil {
				return err
			}

			dataBody := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataBody)
			if err != nil {
				return errors.New(handler + " —— 格式错误")
			}

			if dataBody["type"] == "file" {
				var file model.ModeDataFileStruct
				file.Key = dataBody["key"]
				file.Src = dataBody["src"]
				file.Description = dataBody["desc"]
				file.Type = "file"
				router.Files = append(router.Files, file)

			} else {
				var text model.ModeDataTextStruct
				text.Key = dataBody["key"]
				text.Value = dataBody["value"]
				text.Description = dataBody["desc"]
				text.Type = "text"
				router.Texts = append(router.Texts, text)
			}
		}

	} else if strings.Contains(handler, "@pmHeader") {
		tmpHeader := strings.Split(handler, "@pmHeader")
		if len(tmpHeader) > 1 {
			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpHeader[1]))
			if err != nil {
				return err
			}

			dataHeader := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataHeader)
			if err != nil {
				return errors.New(handler + " —— 格式错误")
			}

			var header model.HeaderStruct
			header.Key = dataHeader["key"]
			header.Name = dataHeader["key"]
			header.Type = "text"
			header.Description = dataHeader["desc"]
			header.Value = dataHeader["value"]
			router.Headers = append(router.Headers, header)
		}
	}
	return nil
}
