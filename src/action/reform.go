package action

import (
	"bytes"
	"errors"
	"io/ioutil"
)

// FindRouters 查找路由用
type FindRouters struct {
	RootRouter   bool
	ParentName   string
	VariableName string
	Path         string
	Handlers     []FindHandlers
	SubRouter    []FindRouters
}

// FindHandlers 查找控制器函数
type FindHandlers struct {
	HandlerPackageName string
	HandlerName        string
}

// RouterGroups 路由组 缓存
var RouterGroups []FindRouters

// ReformFile 逐行遍历，添加 Api 文件
func ReformFile(rootPath string, ignoreFolders []string) error {
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
			var finalFile [][]byte
			// 是否需要写文件
			mark := false

			filePath := rootPath + "/" + file.Name()
			body, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			bodySlice := bytes.Split(body, []byte("\n"))

			for k, b := range bodySlice {
				if k != 0 {
					if !bytes.Contains(bodySlice[k-1], []byte("@ApiQuery(")) {
						appendAPIs(&finalFile, b, &mark)
					}
				} else {
					appendAPIs(&finalFile, b, &mark)
				}
				finalFile = append(finalFile, b)
			}

			// 写文件
			if mark {

				fileContent := bytes.Join(finalFile, []byte("\n"))
				err = WriteFiles(filePath, fileContent)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// appendAPIs appendAPIs
func appendAPIs(finalFile *[][]byte, src []byte, mark *bool) {
	if bytes.Contains(src, []byte(".QueryValue(\"")) && bytes.Contains(src, []byte("\")")) && !bytes.Contains(src, []byte("//")) {
		i := bytes.Index(src, []byte(".QueryValue(\""))
		f := bytes.Index(src, []byte("\")"))
		key := []byte(src[i+13 : f])

		var queryBytes [][]byte
		queryBytes = append(queryBytes, []byte("// @ApiQuery(key=\""))
		queryBytes = append(queryBytes, bytes.TrimSpace([]byte(key)))
		queryBytes = append(queryBytes, []byte("\",desc= \"\", value=\"\")"))
		query := bytes.Join(queryBytes, nil)

		*finalFile = append(*finalFile, query)
		*mark = true

	} else if bytes.Contains(src, []byte(".FormValue(\"")) && bytes.Contains(src, []byte("\")")) && !bytes.Contains(src, []byte("//")) {
		i := bytes.Index(src, []byte(".FormValue(\""))
		f := bytes.Index(src, []byte("\")"))
		key := []byte(src[i+12 : f])

		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte("// @ApiBody(key=\""))
		bodyBytes = append(bodyBytes, bytes.TrimSpace([]byte(key)))
		bodyBytes = append(bodyBytes, []byte("\",desc=\"\", value=\"\")"))
		body := bytes.Join(bodyBytes, nil)

		*finalFile = append(*finalFile, body)
		*mark = true
	} else if bytes.Contains(src, []byte(".Router.GROUP(\"")) && !bytes.Contains(src, []byte("//")) {
		// 获取 变量名
		variableSlice := bytes.SplitN(src, []byte(":="), 2)
		if len(variableSlice) != 2 {
			return
		}
		variable := bytes.TrimSpace(variableSlice[0])

		pathRaw := variableSlice[1]
		i := bytes.Index(pathRaw, []byte("\""))
		pathRaw = bytes.Replace(pathRaw, []byte("\""), []byte("|"), 1)
		f := bytes.Index(pathRaw, []byte("\""))

		var rootRouter FindRouters
		rootRouter.VariableName = string(variable)
		rootRouter.Path = string(pathRaw[i+1 : f])
		rootRouter.RootRouter = true

	} else if bytes.Contains(src, []byte(".GROUP(\"")) && !bytes.Contains(src, []byte("//")) {

	} else if bytes.Contains(src, []byte(".GET(\"")) && !bytes.Contains(src, []byte("//")) {

	} else if bytes.Contains(src, []byte(".POST(\"")) && !bytes.Contains(src, []byte("//")) {

	}
}

func parseRouterProperties(src []byte, isRootRouter bool) (FindRouters, error) {
	variableSlice := bytes.SplitN(src, []byte(":="), 2)
	if len(variableSlice) != 2 {
		return FindRouters{}, errors.New("wrong format")
	}
	variable := bytes.TrimSpace(variableSlice[0])

	pathRaw := variableSlice[1]
	i := bytes.Index(pathRaw, []byte("\""))
	pathRaw = bytes.Replace(pathRaw, []byte("\""), []byte("|"), 1)
	f := bytes.Index(pathRaw, []byte("\""))

	var rootRouter FindRouters
	rootRouter.VariableName = string(variable)
	rootRouter.Path = string(pathRaw[i+1 : f])
	rootRouter.RootRouter = isRootRouter

	// 路由函数
	// TODO
	// handlerSlice := bytes.Split(variableSlice[1], []byte(","))

	return FindRouters{}, nil
}
