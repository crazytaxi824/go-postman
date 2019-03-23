package action

import (
	"io/ioutil"
	"log"
	"strings"
)

// FindRouters 查找路由用
type FindRouters struct {
	ParentName   string
	VariableName string
	Path         string
	// Handlers     []FindHandlers
	HandlersName []string
	Method       string
	SubRouter    []FindRouters
}

// FindHandlers 查找控制器函数
type FindHandlers struct {
	HandlerPackageName string
	HandlerName        string
}

// routerGroups 路由组 缓存
var routerGroups map[string]FindRouters

// ReformFile 逐行遍历，添加 Api 文件
func ReformFile(rootPath string, ignoreFolders []string) error {
	routerGroups = make(map[string]FindRouters)

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if !inSlice(file.Name(), ignoreFolders) {
				ReformFile(rootPath+"/"+file.Name(), ignoreFolders)
			}
		} else {
			var finalFile []string
			// 是否需要写文件
			mark := false

			filePath := rootPath + "/" + file.Name()
			body, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			bodySlice := strings.Split(string(body), "\n")

			for k, str := range bodySlice {
				if k != 0 {
					if !strings.Contains(bodySlice[k-1], "@Api") {
						err = appendAPIs(&finalFile, str, &mark)
						if err != nil {
							continue
						}
					}
				} else {
					err := appendAPIs(&finalFile, str, &mark)
					if err != nil {
						continue
					}
				}
				finalFile = append(finalFile, str)
			}

			// TODO 分析 FindRouter

			// 写文件
			if mark {
				// for _, v := range routerGroups {
				// 	log.Println(v)
				// }

				fileContent := strings.Join(finalFile, "\n")
				log.Println(fileContent)
				// err = WriteFiles([]byte(filePath), fileContent)
				// if err != nil {
				// 	return err
				// }
			}
		}
	}
	return nil
}

// appendAPIs appendAPIs
func appendAPIs(finalFile *[]string, src string, mark *bool) (err error) {
	if strings.Contains(src, ".QueryValue(\"") && strings.Contains(src, "\")") && !strings.Contains(src, "//") {
		i := strings.Index(src, ".QueryValue(\"")
		f := strings.Index(src, "\")")
		key := strings.TrimSpace(src[i+13 : f])

		query := "// @ApiQuery(key=\"" + key + "\", desc= \"\", value=\"\")"

		*finalFile = append(*finalFile, query)
		*mark = true

	} else if strings.Contains(src, ".FormValue(\"") && strings.Contains(src, "\")") && !strings.Contains(src, "//") {
		i := strings.Index(src, ".FormValue(\"")
		f := strings.Index(src, "\")")
		key := strings.TrimSpace(src[i+12 : f])

		body := "// @ApiBody(key=\"" + key + "\", desc=\"\", value=\"\")"

		*finalFile = append(*finalFile, body)
		*mark = true

	} else if strings.Contains(src, ".GROUP(\"") && !strings.Contains(src, "//") {

		// *mark = true
	} else if strings.Contains(src, ".GET(\"") && !strings.Contains(src, "//") {

		// *mark = true
	} else if strings.Contains(src, ".POST(\"") && !strings.Contains(src, "//") {

		// *mark = true
	}
	return nil
}

func parseRouterProperties(src string) {

}
