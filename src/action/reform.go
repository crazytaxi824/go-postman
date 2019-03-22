package action

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"path"
)

// FindRouters 查找路由用
type FindRouters struct {
	RootRouter   bool
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
var routerGroups []FindRouters

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
					if !bytes.Contains(bodySlice[k-1], []byte("@Api")) {
						err = appendAPIs(&finalFile, b, &mark)
						if err != nil {
							continue
						}
					}
				} else {
					err := appendAPIs(&finalFile, b, &mark)
					if err != nil {
						continue
					}
				}
				finalFile = append(finalFile, b)
			}

			// TODO 分析 FindRouter

			// 写文件
			if mark {
				for _, v := range routerGroups {
					log.Println(v)
				}

				fileContent := bytes.Join(finalFile, []byte("\n"))
				log.Println(string(fileContent))
				// err = WriteFiles(filePath, fileContent)
				// if err != nil {
				// 	return err
				// }
			}
		}
	}
	return nil
}

// appendAPIs appendAPIs
func appendAPIs(finalFile *[][]byte, src []byte, mark *bool) (err error) {
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

	} else if bytes.Contains(src, []byte(".GROUP(\"")) && !bytes.Contains(src, []byte("//")) {
		if bytes.Contains(src, []byte(".Router.")) {
			err = parseGroupRouterProperties(src, true, "")
			if err != nil {
				return err
			}

		} else {
			err = parseGroupRouterProperties(src, false, "")
			if err != nil {
				return err
			}
		}
		*mark = true

	} else if bytes.Contains(src, []byte(".GET(\"")) && !bytes.Contains(src, []byte("//")) {
		if bytes.Contains(src, []byte(".Router.")) {
			r, err := parseRootProperties(src, "GET")
			if err != nil {
				log.Println(err.Error())
				return err
			}

			var queryBytes [][]byte
			queryBytes = append(queryBytes, []byte("// @ApiQuery(name=\""))
			queryBytes = append(queryBytes, bytes.TrimSpace([]byte(path.Base(r.Path)+" "+r.ParentName)))
			queryBytes = append(queryBytes, []byte("\", path= \""))
			queryBytes = append(queryBytes, []byte(r.Path))
			queryBytes = append(queryBytes, []byte("\", method=\""))
			queryBytes = append(queryBytes, []byte(r.Method))
			queryBytes = append(queryBytes, []byte("\")"))
			query := bytes.Join(queryBytes, nil)

			*finalFile = append(*finalFile, query)
		} else {
			err := parseRouterProperties(src, "GET")
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
		*mark = true
	} else if bytes.Contains(src, []byte(".POST(\"")) && !bytes.Contains(src, []byte("//")) {
		if bytes.Contains(src, []byte(".Router.")) {
			// r, err := parseRootProperties(src, "POST")
			// if err != nil {
			// 	log.Println(err.Error())
			// 	return err
			// }
		} else {
			err := parseRouterProperties(src, "POST")
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
		*mark = true
	}
	return nil
}

func parseRootProperties(src []byte, method string) (*FindRouters, error) {
	pathRaw := src
	i := bytes.Index(pathRaw, []byte("\""))
	pathRaw = bytes.Replace(pathRaw, []byte("\""), []byte("|"), 1)
	f := bytes.Index(pathRaw, []byte("\""))

	var router FindRouters
	router.Path = string(bytes.TrimSpace(pathRaw[i+1 : f]))
	router.RootRouter = true
	router.Method = method

	// 路由处理器
	// TODO 分析handlers
	handlerRawSlice := bytes.Split(src, []byte(")"))
	if len(handlerRawSlice) != 2 {
		return nil, errors.New("no handler")
	}
	handlerSlice := bytes.Split(handlerRawSlice[0], []byte(","))
	lenHandler := len(handlerSlice)
	if lenHandler < 2 {
		return nil, errors.New("no handler")
	}

	for i := 1; i < lenHandler; i++ {
		router.HandlersName = append(router.HandlersName, string(bytes.TrimSpace(handlerSlice[i])))
	}

	return &router, nil
}

func parseGroupRouterProperties(src []byte, isRootRouter bool, method string) error {
	variableSlice := bytes.SplitN(src, []byte(":="), 2)
	if len(variableSlice) != 2 {
		return errors.New("wrong format")
	}

	pathRaw := variableSlice[1]
	i := bytes.Index(pathRaw, []byte("\""))
	pathRaw = bytes.Replace(pathRaw, []byte("\""), []byte("|"), 1)
	f := bytes.Index(pathRaw, []byte("\""))

	var router FindRouters
	router.VariableName = string(bytes.TrimSpace(variableSlice[0]))
	router.Path = string(bytes.TrimSpace(pathRaw[i+1 : f]))
	router.RootRouter = isRootRouter
	router.Method = method

	// 路由处理器
	// TODO 分析handlers
	handlerRawSlice := bytes.Split(variableSlice[1], []byte(")"))
	if len(handlerRawSlice) != 2 {
		return errors.New("no handler")
	}
	handlerSlice := bytes.Split(handlerRawSlice[0], []byte(","))
	lenHandler := len(handlerSlice)
	if lenHandler < 2 {
		return errors.New("no handler")
	}

	for i := 1; i < lenHandler; i++ {
		router.HandlersName = append(router.HandlersName, string(bytes.TrimSpace(handlerSlice[i])))
	}

	// 查找parentName
	if !isRootRouter {
		pNameSlice := bytes.Split(variableSlice[1], []byte("("))
		if len(pNameSlice) < 2 {
			return errors.New("wrong format")
		}

		pNameSlice = bytes.Split(pNameSlice[0], []byte("."))
		if len(pNameSlice) < 2 {
			return errors.New("wrong format")
		}

		lenPM := len(pNameSlice)
		pName := pNameSlice[:lenPM-1]

		router.ParentName = string(bytes.TrimSpace(bytes.Join(pName, []byte("."))))

		// 查找上级router
		for k := range routerGroups {
			routerGroups[k].findingParentRouter(router)
		}
		return nil
	}

	routerGroups = append(routerGroups, router)

	return nil
}

func (router *FindRouters) findingParentRouter(r FindRouters) {
	if router.VariableName == r.ParentName {
		router.SubRouter = append(router.SubRouter, r)
	} else {
		for k := range router.SubRouter {
			router.SubRouter[k].findingParentRouter(r)
		}
	}
}

func parseRouterProperties(src []byte, method string) error {
	pathRaw := src
	i := bytes.Index(pathRaw, []byte("\""))
	pathRaw = bytes.Replace(pathRaw, []byte("\""), []byte("|"), 1)
	f := bytes.Index(pathRaw, []byte("\""))

	var router FindRouters
	router.Path = string(bytes.TrimSpace(pathRaw[i+1 : f]))
	router.RootRouter = false
	router.Method = method

	// 路由处理器
	// TODO 分析handlers
	handlerRawSlice := bytes.Split(src, []byte(")"))
	if len(handlerRawSlice) != 2 {
		return errors.New("no handler")
	}

	handlerSlice := bytes.Split(handlerRawSlice[0], []byte(","))
	lenHandler := len(handlerSlice)
	if lenHandler < 2 {
		return errors.New("no handler")
	}

	for i := 1; i < lenHandler; i++ {
		router.HandlersName = append(router.HandlersName, string(bytes.TrimSpace(handlerSlice[i])))
	}

	// 查找parentName
	pNameSlice := bytes.Split(src, []byte("("))
	if len(pNameSlice) < 2 {
		return errors.New("wrong format")
	}

	pNameSlice = bytes.Split(pNameSlice[0], []byte("."))
	if len(pNameSlice) < 2 {
		return errors.New("wrong format")
	}

	lenPM := len(pNameSlice)
	pName := pNameSlice[:lenPM-1]

	router.VariableName = string(bytes.TrimSpace(bytes.Join(pName, []byte("."))))

	// 查找自己的router
	for k := range routerGroups {
		routerGroups[k].findingSelfRouter(router)
	}

	return nil
}

func (router *FindRouters) findingSelfRouter(r FindRouters) {
	if len(router.SubRouter) < 1 && router.VariableName == r.VariableName {
		router.Path = router.Path + r.Path
		router.Method = r.Method
		for _, v := range r.HandlersName {
			router.HandlersName = append(router.HandlersName, v)
		}
	} else {
		for k := range router.SubRouter {
			router.SubRouter[k].findingSelfRouter(r)
		}
	}
}
