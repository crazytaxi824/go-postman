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
	ignoreFolders := flag.String("i", "vendor", "folder names, ignore multi folders, using | to split")
	outputPath := flag.String("o", "./newPostman.json", "output file name")
	specify := flag.String("s", "", "specify file suffix, eg: .go")
	format := flag.Bool("format", false, "write API to your files, package HttpDispatch only")
	flag.Parse()

	// *ignoreFile = "vendor|action|model|main"
	// *format = true
	// *specify = ".go"

	ignoreFolderSlice := strings.Split(*ignoreFolders, "|")
	for k := range ignoreFolderSlice {
		ignoreFolderSlice[k] = strings.TrimSpace(ignoreFolderSlice[k])
	}

	action.HandlerMap = make(map[string][]action.RawHandlerStruct)
	action.ProjectFiles = make(map[string][]action.AllFiles)

	if *format {

		err = action.ReformFile(*rootPath, ignoreFolderSlice, *specify)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = action.AnalysisFindRouter()
		if err != nil {
			log.Println(err.Error())
			return
		}

		for _, files := range action.ProjectFiles {
			for _, file := range files {
				if file.FormatMark {
					log.Println("file formated: " + file.FileName)

					// 写文件
					err = action.WriteFiles(file.FileName, []byte(file.Content))
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}

		return
	}

	// 读取文件夹下所有go文件 -----------------------------------------------------
	var serverPath string
	var routers []action.RawRouterStruct
	err = action.ReadAllFiles(*rootPath, &serverPath, ignoreFolderSlice, &routers, *specify)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 把处理器中的所有参数传入路由 ------------------------------------------------
	for k := range routers {
		routers[k].HandlersToRouters()
	}

	// 生成url，生成 header --------------------------------------------------
	for k := range routers {
		err = routers[k].GenHeaderAndURLStruct(serverPath)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	// 给路由分组 -----------------------------------------------------------
	groups := action.GroupRouters(routers)

	if len(groups) == 0 {
		log.Println("no @Api detected!")
		return
	}

	// 生成 body ------------------------------------------------------------
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
