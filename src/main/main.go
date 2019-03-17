package main

import (
	"action"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"model"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// JSON JSON
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	log.SetFlags(log.Lshortfile)

	rootPath := "./src"

	// 读取文件夹下所有go文件 -----------------------------------------------------
	var serverPath string
	var routers []action.RawRouterStruct
	var rawHandlerSlice []string
	err := readAllFiles(rootPath, &serverPath, &routers, &rawHandlerSlice)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(" ----------------------------------------- ")

	// 把处理器传入路由 ---------------------------------------------------------
	action.HandlersToRouters(&routers, rawHandlerSlice)

	// 生成url，生成 header --------------------------------------------------
	action.GenHeaderAndURLStruct(&routers, serverPath)

	// // 给路由分组
	groups := action.GroupRouters(routers)

	// ---------------------------------------------------------------------
	// 生成 body

	action.GenGroupBody(&groups)

	// 生成 folder & item -----------------------------------------------------
	folders, items := action.GenFolderAndItemStruct(groups)

	// 生成 pm 文件 ---------------------------------------------
	var pm model.PostmanStruct
	pm.Info.Name = "postman-test"
	pm.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	for k := range folders {
		pm.Item = append(pm.Item, folders[k])
	}

	for k := range items {
		pm.Item = append(pm.Item, items[k])
	}

	// 生成 pm 文件
	// 生成代码不含 \u0026
	bf := bytes.NewBuffer([]byte{})
	jsonEncode := JSON.NewEncoder(bf)
	jsonEncode.SetEscapeHTML(false)
	err = jsonEncode.Encode(pm)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(bf.String())

	// TODO 输出文件

}

// readAllFiles 读取所有数据
func readAllFiles(rootPath string, serverPath *string, routers *[]action.RawRouterStruct, rawHandlerSlice *[]string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		if file.IsDir() {
			if file.Name() != "vendor" {
				readAllFiles(rootPath+"/"+file.Name(), serverPath, routers, rawHandlerSlice)
			}
		} else {
			if path.Ext(path.Base(file.Name())) == ".go" {

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
									var router action.RawRouterStruct
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
	}
	return nil
}
