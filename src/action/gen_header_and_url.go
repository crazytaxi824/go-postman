package action

import (
	"errors"
	"log"
	"model"
	"strings"
)

// GenHeaderAndURLStruct 生成url和header
func (router *RawRouterStruct) GenHeaderAndURLStruct(serverPath string) (err error) {
	err = router.rawURLtoRequestURL(serverPath + router.RouterPath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// 生成 header

	var hasContentType bool
	if len(router.Files) == 0 && len(router.Texts) > 0 {
		// 判断header中是否存在"Content-Type"
		for _, v := range router.Headers {
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

			router.Headers = append(router.Headers, headerUrlencode)
		}
	} else {
		if len(router.Headers) < 1 {
			router.Headers = make([]model.HeaderStruct, 0)
		}
	}
	return nil
}

// rawURLtoRequestURL 通过query生成为 model.URLStruct
func (router *RawRouterStruct) rawURLtoRequestURL(path string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var urlStruct model.URLStruct

	prefix := strings.Split(path, "://")
	if len(prefix) < 2 {
		return errors.New("server path error")
	}

	// 生成 protocol
	urlStruct.Protocol = prefix[0]

	// 生成path
	rawPath := strings.Split(prefix[1], "/")

	// 生成 host and port
	rawHostAndPort := strings.Split(rawPath[0], ":")
	length := len(rawHostAndPort)
	if length == 2 {
		urlStruct.Port = rawHostAndPort[1]
	} else if length > 2 {
		return errors.New("server path error")
	}
	rawHost := strings.Split(rawHostAndPort[0], ".")
	urlStruct.Host = rawHost

	// 生成 path
	// 判断rawPath长度
	if len(rawPath) > 1 {
		urlStruct.Path = rawPath[1:]
	}

	// 生成 query
	urlStruct.Query = router.Querys

	// 生成 raw
	urlStruct.Raw = path
	if len(router.Querys) > 0 {
		for _, v := range router.Querys {
			urlStruct.Raw = urlStruct.Raw + "&" + v.Key + "=" + v.Value
		}
		urlStruct.Raw = strings.Replace(urlStruct.Raw, "&", "?", 1)
	}

	router.URL = urlStruct

	return nil
}
