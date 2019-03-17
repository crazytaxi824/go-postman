package action

import (
	"bytes"
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
				trimBody := string(bytes.TrimSpace(v))
				if len(trimBody) > 4 {
					if trimBody[:2] == "//" {
						if strings.Contains(trimBody, "@pmServer") {

							// 处理 serverPath
							tmp := strings.Split(trimBody, "@pmServer")
							if len(tmp) > 1 {
								data := make(map[string]string)
								err = JSON.UnmarshalFromString(tmp[1], &data)
								if err != nil {
									// log.Println(string(v) + "pmServer 格式错误")
									continue
								}
								// if _, ok := data["path"]; !ok {
								// 	log.Println("pmServer 没有\"path\"参数")
								// }
								*serverPath = data["path"]
							}
						} else if strings.Contains(trimBody, "@pmRouter") {

							// 处理 router
							tmp := strings.Split(trimBody, "@pmRouter")
							if len(tmp) > 1 {
								var router RawRouterStruct
								err = JSON.UnmarshalFromString(tmp[1], &router)
								router.Method = strings.ToUpper(router.Method)
								if err != nil {
									// log.Println(string(v) + " —— 格式错误")
									continue
								}
								*routers = append(*routers, router)
							}
						} else if strings.Contains(trimBody, "@pmHandler") || strings.Contains(trimBody, "@pmQuery") || strings.Contains(trimBody, "@pmBody") || strings.Contains(trimBody, "@pmHeader") {
							// 处理 pmHandler, pmQuery, pmBody, pmHeader
							*rawHandlerSlice = append(*rawHandlerSlice, trimBody)
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
