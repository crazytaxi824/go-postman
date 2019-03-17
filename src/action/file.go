package action

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// ReadAllFiles 读取所有数据
func ReadAllFiles(rootPath string, serverPath *string, ignoreFolders []string, routers *[]RawRouterStruct, rawHandlerSlice *[]string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		if file.IsDir() {
			// if file.Name() != "vendor" {
			if !inSlice(file.Name(), ignoreFolders) {
				ReadAllFiles(rootPath+"/"+file.Name(), serverPath, ignoreFolders, routers, rawHandlerSlice)
			}
		} else {
			body, err := ioutil.ReadFile(rootPath + "/" + file.Name())
			if err != nil {
				return err
			}

			bodySlice := bytes.Split(body, []byte("\n"))

			for _, v := range bodySlice {

				res := string(bytes.TrimSpace(v))

				if len(res) > 4 {
					if res[:2] == "//" {
						if strings.Contains(res, "@pmServer") {
							// 处理 serverPath
							tmp := strings.Split(res, "@pmServer")
							if len(tmp) > 1 {
								data := make(map[string]string)
								err = JSON.UnmarshalFromString(tmp[1], &data)
								if err != nil {
									continue
									// return errors.New(res + " —— 格式错误")
								}
								if _, ok := data["path"]; !ok {
									return errors.New(res + " —— 没有 path 参数")
								}
								*serverPath = data["path"]
							}
						} else if strings.Contains(res, "@pmRouter") {
							// 处理 router
							tmp := strings.Split(res, "@pmRouter")
							if len(tmp) > 1 {
								var router RawRouterStruct
								err = JSON.UnmarshalFromString(tmp[1], &router)
								router.Method = strings.ToUpper(router.Method)
								if err != nil {
									continue
									// return errors.New(res + " —— 格式错误")
								}
								*routers = append(*routers, router)
							}
						} else if strings.Contains(res, "@pmHandler") || strings.Contains(res, "@pmQuery") || strings.Contains(res, "@pmBody") || strings.Contains(res, "@pmHeader") {
							*rawHandlerSlice = append(*rawHandlerSlice, res)
						}
					}
				}
			}
		}
	}
	return nil
}

// WriteFiles 写文件
func WriteFiles(filename string, content []byte) error {
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func inSlice(s string, ss []string) bool {
	for k := range ss {
		if ss[k] == s {
			return true
		}
	}
	return false
}
