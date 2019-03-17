package main

import (
	"action"
	"bytes"
	"flag"
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
	var err error

	rootPath := flag.String("p", "./src", "read all files from this path, including sub folders")
	ignoreFile := flag.String("i", "vendor", "folder names, ignore multi folders, using | to split")
	outputPath := flag.String("o", "./newPostman.json", "output file name")
	flag.Parse()

	// rootPath := "./src"
	ignoreFiles := strings.Split(*ignoreFile, "|")
	for k := range ignoreFiles {
		ignoreFiles[k] = strings.TrimSpace(ignoreFiles[k])
	}

	// 读取文件夹下所有go文件 -----------------------------------------------------
	var serverPath string
	var routers []action.RawRouterStruct
	var rawHandlerSlice []string
	err = action.ReadAllFiles(*rootPath, &serverPath, ignoreFiles, &routers, &rawHandlerSlice)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 把处理器传入路由 ---------------------------------------------------------
	for k := range routers {
		routers[k].HandlersToRouters(rawHandlerSlice)
	}

	// 生成url，生成 header --------------------------------------------------
	for k := range routers {
		err = routers[k].GenHeaderAndURLStruct(serverPath)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	// 给路由分组
	groups := action.GroupRouters(routers)

	// ---------------------------------------------------------------------
	// 生成 body

	action.GenGroupBody(&groups)

	// 生成 folder & item -----------------------------------------------------
	folders, items := action.GenFolderAndItemStruct(groups)

	// 生成 pm 文件 ---------------------------------------------
	var pm model.PostmanStruct
	outputFileSuffix := path.Ext(*outputPath)
	outputFileName := strings.TrimSuffix(path.Base(*outputPath), outputFileSuffix)
	pm.Info.Name = outputFileName
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

	// 输出文件
	err = action.WriteFiles(*outputPath, bf.Bytes())
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("gen Postman file completed!")
}
