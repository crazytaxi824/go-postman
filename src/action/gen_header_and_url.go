package action

import (
	"errors"
	"log"
	"model"
	"strings"
)

// GenHeaderAndURLStruct 生成url和header
func GenHeaderAndURLStruct(routersPointer *[]RawRouterStruct, serverPath string) (err error) {
	routers := *routersPointer
	for k := range routers {
		routers[k].URL, err = rawURLtoRequestURL(serverPath+routers[k].RouterPath, routers[k].Querys)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		// 生成 header

		var hasContentType bool
		if len(routers[k].Files) == 0 && len(routers[k].Texts) > 0 {
			// 判断header中是否存在"Content-Type"
			for _, v := range routers[k].Headers {
				if v.Key == "Content-Type" {
					hasContentType = true
				}
			}

			if !hasContentType {
				var headerUrlencode model.HeaderStruct
				headerUrlencode.Type = "text"
				headerUrlencode.Key = "Content-Type"
				headerUrlencode.Name = "Content-Type"
				headerUrlencode.Value = "application/x-www-form-urlencoded"
				headerUrlencode.Description = ""

				routers[k].Headers = append(routers[k].Headers, headerUrlencode)
			}
		} else {
			if len(routers[k].Headers) < 1 {
				routers[k].Headers = make([]model.HeaderStruct, 0)
			}
		}
	}
	return nil
}

// rawURLtoRequestURL 通过query生成为 model.URLStruct
func rawURLtoRequestURL(path string, query []model.QueryStruct) (model.URLStruct, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var urlStruct model.URLStruct

	prefix := strings.Split(path, "://")
	if len(prefix) < 2 {
		return model.URLStruct{}, errors.New("请求路径错误")
	}

	// 生成 protocol
	urlStruct.Protocol = prefix[0]

	rawPath := strings.Split(prefix[1], "/")

	// 生成 host and port
	rawHostAndPort := strings.Split(rawPath[0], ":")
	length := len(rawHostAndPort)
	if length == 2 {
		urlStruct.Port = rawHostAndPort[1]
	} else if length > 2 {
		return model.URLStruct{}, errors.New("请求路径错误")
	}
	rawHost := strings.Split(rawHostAndPort[0], ".")
	urlStruct.Host = rawHost

	// 生成 path
	// 判断rawPath长度
	if len(rawPath) > 1 {
		urlStruct.Path = rawPath[1:]
	}

	// 生成 query
	urlStruct.Query = query

	// 生成 raw
	urlStruct.Raw = path
	if len(query) > 0 {
		for _, v := range query {
			urlStruct.Raw = urlStruct.Raw + "&" + v.Key + "=" + v.Value
		}
		urlStruct.Raw = strings.Replace(urlStruct.Raw, "&", "?", 1)
	}

	return urlStruct, nil
}
