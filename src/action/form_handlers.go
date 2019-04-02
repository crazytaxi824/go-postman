package action

import (
	"errors"
	"log"
	"model"
	"strings"
)

// SaveHandlers 将handler和对应的body放入 handlerMap 缓存中
func SaveHandlers(handlers []string) {
	handlerMap = make(map[string][]RawHandlerStruct)

	var handlerIndex []int
	for k, h := range handlers {
		if strings.Contains(h, "@ApiHandler") {
			handlerIndex = append(handlerIndex, k)
		}
	}

	// 获取 handler name
	for i := 0; i < len(handlerIndex); i++ {
		tmp := strings.Split(handlers[handlerIndex[i]], "@ApiHandler")
		if len(tmp) > 1 {
			handlerRef, err := ParsePMstructToJSONformat(strings.TrimSpace(tmp[1]))
			if err != nil {
				log.Println("warning: format error ——" + handlers[handlerIndex[i]])
				continue
			}

			data := make(map[string]string)
			err = JSON.UnmarshalFromString(handlerRef, &data)
			if err != nil {
				log.Println("warning: format error —— " + handlers[handlerIndex[i]])
				continue
			}

			// // 匹配 handler name 和 router name
			// // 如果 handler name 不匹配则会在这里被丢弃

			if i > len(handlerIndex)-2 {
				for _, handler := range handlers[handlerIndex[i]+1:] {
					// 将参数传入对应的 handlerMap 中
					err = passQueryBodyHeaderToHandler(handler, data["name"])
					if err != nil {
						log.Println("warning: format error —— " + handler)
						log.Println(err.Error())
						return
					}
				}
			} else {
				for _, handler := range handlers[handlerIndex[i]+1 : handlerIndex[i+1]] {
					// 将参数传入对应的 handlerMap 中
					err = passQueryBodyHeaderToHandler(handler, data["name"])
					if err != nil {
						log.Println("warning: format error —— " + handler)
						log.Println(err.Error())
						return
					}
				}
			}
		}
	}

	return
}

// 将 body 传入 handler 缓存
func passQueryBodyHeaderToHandler(handler, handlerName string) error {
	var handlerStruct RawHandlerStruct

	var queryKeysName []string
	var bodyKeysName []string
	var headerKeysName []string

	if strings.Contains(handler, "@ApiQuery") {
		tmpQuery := strings.Split(handler, "@ApiQuery")
		if len(tmpQuery) > 1 {

			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpQuery[1]))
			if err != nil {
				return err
			}

			dataQuery := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataQuery)
			if err != nil {
				return errors.New(handler + " —— format error")
			}

			var query model.QueryStruct
			query.Key = dataQuery["key"]
			query.Value = dataQuery["value"]
			query.Description = dataQuery["desc"]

			// 检查key是否重复
			if inSlice(query.Key, queryKeysName) {
				return errors.New("duplicate QUERY key —— router: " + handlerName + ", key: " + query.Key)
			}

			handlerStruct.Querys = append(handlerStruct.Querys, query)
			queryKeysName = append(queryKeysName, query.Key)
		}
	} else if strings.Contains(handler, "@ApiBody") {
		tmpBody := strings.Split(handler, "@ApiBody")
		if len(tmpBody) > 1 {
			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpBody[1]))
			if err != nil {
				return err
			}

			dataBody := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataBody)
			if err != nil {
				return errors.New(handler + " —— format error")
			}

			if dataBody["type"] == "file" {
				var file model.ModeDataFileStruct
				file.Key = dataBody["key"]
				file.Src = dataBody["src"]
				file.Description = dataBody["desc"]
				file.Type = "file"

				if inSlice(file.Key, bodyKeysName) {
					return errors.New("duplicate BODY key —— router: " + handlerName + ", key: " + file.Key)
				}

				handlerStruct.Files = append(handlerStruct.Files, file)
				bodyKeysName = append(bodyKeysName, file.Key)

			} else {
				var text model.ModeDataTextStruct
				text.Key = dataBody["key"]
				text.Value = dataBody["value"]
				text.Description = dataBody["desc"]
				text.Type = "text"

				// 检查key是否重复
				if inSlice(text.Key, bodyKeysName) {
					return errors.New("duplicate BODY key —— router: " + handlerName + ", key: " + text.Key)
				}

				handlerStruct.Texts = append(handlerStruct.Texts, text)
				bodyKeysName = append(bodyKeysName, text.Key)
			}
		}

	} else if strings.Contains(handler, "@ApiHeader") {
		tmpHeader := strings.Split(handler, "@ApiHeader")
		if len(tmpHeader) > 1 {
			ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmpHeader[1]))
			if err != nil {
				return err
			}

			dataHeader := make(map[string]string)
			err = JSON.UnmarshalFromString(ref, &dataHeader)
			if err != nil {
				return errors.New(handler + " —— format error")
			}

			var header model.HeaderStruct
			header.Key = dataHeader["key"]
			header.Name = dataHeader["key"]
			header.Type = "text"
			header.Description = dataHeader["desc"]
			header.Value = dataHeader["value"]

			// 检查key是否重复
			if inSlice(header.Key, headerKeysName) {
				return errors.New("duplicate BODY key —— router: " + handlerName + ", key: " + header.Key)
			}

			handlerStruct.Headers = append(handlerStruct.Headers, header)
			headerKeysName = append(headerKeysName, header.Key)
		}
	}

	handlerMap[handlerName] = append(handlerMap[handlerName], handlerStruct)
	return nil
}
