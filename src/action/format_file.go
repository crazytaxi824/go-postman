package action

import (
	"errors"
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
	// RootRouter   bool
	HandlersName []string
	Method       string
}

// FindHandlers 查找控制器函数
type FindHandlers struct {
	HandlerPackageName string
	HandlerName        string
}

// routerGroups 路由组 缓存
var routerGroups map[string]FindRouters
var rootRouterGroups []FindRouters

// ReformFile 逐行遍历，添加 Api 文件
func ReformFile(rootPath string, ignoreFolders []string) error {
	routerGroups = make(map[string]FindRouters)

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	var fileName string
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
				// if k != 0 {
				// 	if !strings.Contains(bodySlice[k-1], "@Api") {
				// 		err = appendAPIs(&finalFile, str, &mark)
				// 		if err != nil {
				// 			continue
				// 		}
				// 	}
				// } else {
				// 	err := appendAPIs(&finalFile, str, &mark)
				// 	if err != nil {
				// 		continue
				// 	}
				// }

				apiStr, err := appendAPIsNew(str)
				if err != nil {
					finalFile = append(finalFile, str)
					continue
				}

				if apiStr != "" {
					if k > 0 {
						if strings.Contains(bodySlice[k-1], "@Api") && strings.Contains(bodySlice[k-1], "//") && strings.TrimSpace(bodySlice[k-1]) != strings.TrimSpace(apiStr) {
							finalFile = finalFile[:len(finalFile)-1]
							finalFile = append(finalFile, apiStr)
							fileName = file.Name()
							mark = true
						} else if strings.TrimSpace(bodySlice[k-1]) != strings.TrimSpace(apiStr) {
							finalFile = append(finalFile, apiStr)
							fileName = file.Name()
							mark = true
						}
					} else {
						finalFile = append(finalFile, apiStr)
						fileName = file.Name()
						mark = true
					}
				}

				finalFile = append(finalFile, str)
			}

			// TODO 分析 FindRouter

			// 写文件
			if mark {
				log.Println("file formated: " + rootPath + fileName)
				fileContent := strings.Join(finalFile, "\n")

				// 写文件
				err = WriteFiles(filePath, []byte(fileContent))
				if err != nil {
					return err
				}
			}
		}
	}

	// for _, v := range rootRouterGroups {
	// 	log.Println(v)
	// }

	return nil
}

func appendAPIsNew(src string) (apiStr string, err error) {
	if strings.Contains(src, ".QueryValue(\"") && strings.Contains(src, "\")") && !strings.Contains(src, "//") {
		i := strings.Index(src, ".QueryValue(\"")
		f := strings.Index(src, "\")")
		key := strings.TrimSpace(src[i+13 : f])

		query := "// @ApiQuery(key=\"" + key + "\", desc= \"\", value=\"\")"

		return query, nil

	} else if strings.Contains(src, ".FormValue(\"") && strings.Contains(src, "\")") && !strings.Contains(src, "//") {
		i := strings.Index(src, ".FormValue(\"")
		f := strings.Index(src, "\")")
		key := strings.TrimSpace(src[i+12 : f])

		body := "// @ApiBody(key=\"" + key + "\", desc=\"\", value=\"\")"

		return body, nil

	} else if strings.Contains(src, ".GROUP(\"") && !strings.Contains(src, "//") {
		r, err := parseRouterGroupProperties(src)
		if err != nil {
			return "", err
		}

		if r != "" {
			return r, nil
		}
	} else if !strings.Contains(src, "//") {
		var router FindRouters
		if !router.getRouterMethod(src) {
			return "", nil
		}

		r, err := router.genRouterAPI(src)
		if err != nil {
			return "", err
		}

		if r != "" {
			return r, nil
		}
	}
	return "", nil
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
		r, err := parseRouterGroupProperties(src)
		if err != nil {
			return err
		}

		if r != "" {
			*finalFile = append(*finalFile, r)
			*mark = true
		}
	} else if !strings.Contains(src, "//") {
		var router FindRouters
		if !router.getRouterMethod(src) {
			return nil
		}

		r, err := router.genRouterAPI(src)
		if err != nil {
			return err
		}

		if r != "" {
			*finalFile = append(*finalFile, r)
			*mark = true
		}
	}
	return nil
}

func parseRouterGroupProperties(src string) (apiStr string, err error) {
	var router FindRouters
	// var router FindRouters
	// 判断 :=
	if strings.Contains(src, ":=") {
		// 获取 variableName
		variableSlice := strings.SplitN(src, ":=", 2)
		router.VariableName = strings.TrimSpace(variableSlice[0])

		// 获取 path And Handler
		router.getGroupPathAndHandlers(src)

		// 获取 parentName
		router.getGroupParentName(variableSlice[1])

		routerGroups[router.VariableName] = router
	} else {
		b := router.getRouterMethod(src)
		if !b {
			return "", errors.New("format error")
		}

		apiStr, err = router.genRouterAPI(src)
		if err != nil {
			return "", err
		}

	}

	return apiStr, nil
}

func (router *FindRouters) genRouterAPI(src string) (string, error) {
	routerSlice := strings.Split(src, "."+router.Method)
	groupsRaw := routerSlice[0]
	routerRaw := routerSlice[1]

	err := router.getPathAndHandler(strings.TrimSpace(routerRaw))
	if err != nil {
		return "", err
	}

	if !strings.Contains(groupsRaw, ".GROUP") {
		router.getRouterParentName(src)
		// 向上查找 parent
		router.findingParent(router.ParentName)

		// 添加到 rootRouterGroups
		rootRouterGroups = append(rootRouterGroups, *router)

		routerNameSlice := strings.Split(router.Path, "/")
		var routerName string
		lenName := len(routerNameSlice)
		if lenName < 1 {
			return "", errors.New("no path")
		} else if lenName < 2 {
			routerName = routerNameSlice[0]
		} else {
			routerName = strings.TrimSpace(strings.Join(routerNameSlice[lenName-2:], " "))
		}

		apiStr := "// @ApiRouter(name=\"" + routerName + "\", method=\"" + router.Method + "\", path=\"" + router.Path + "\", group=\"" + router.ParentName + "\", handlers=\"" + strings.Join(router.HandlersName, ",") + "\")"

		return apiStr, nil
	}

	var group FindRouters
	groupsSlice := strings.Split(groupsRaw, ".GROUP")
	for k, v := range groupsSlice {
		if k == 0 {
			continue
		}

		err := group.getPathAndHandler(v)
		if err != nil {
			return "", err
		}
	}

	// 查找parentName
	group.getGroupParentName(groupsRaw)

	router.ParentName = group.ParentName
	router.Path = group.Path + router.Path
	router.HandlersName = append(router.HandlersName, group.HandlersName...)

	// 向上查找 parent
	router.findingParent(group.ParentName)

	// 添加到 rootRouterGroups
	rootRouterGroups = append(rootRouterGroups, *router)

	routerNameSlice := strings.Split(router.Path, "/")
	var routerName string
	lenName := len(routerNameSlice)
	if lenName < 1 {
		return "", errors.New("no path")
	} else if lenName < 2 {
		routerName = routerNameSlice[0]
	} else {
		routerName = strings.TrimSpace(strings.Join(routerNameSlice[lenName-2:], " "))
	}
	apiStr := "// @ApiRouter(name=\"" + routerName + "\", method=\"" + router.Method + "\", path=\"" + router.Path + "\", group=\"" + group.ParentName + "\", handlers=\"" + strings.Join(router.HandlersName, ",") + "\")"

	return apiStr, nil
}

func (router *FindRouters) findingParent(parentName string) {
	router.Path = routerGroups[parentName].Path + router.Path
	router.HandlersName = append(router.HandlersName, routerGroups[parentName].HandlersName...)
	if routerGroups[parentName].ParentName != "" {
		router.findingParent(routerGroups[parentName].ParentName)
	}
}

func (router *FindRouters) getGroupParentName(srcWithourVariableName string) {
	if strings.Contains(srcWithourVariableName, ".Router.GROUP(\"") {
		return
	}

	parentNameSlice := strings.Split(srcWithourVariableName, ".GROUP")
	router.ParentName = strings.TrimSpace(parentNameSlice[0])
	return
}

func (router *FindRouters) getRouterParentName(srcWithourVariableName string) {
	if strings.Contains(srcWithourVariableName, ".Router."+router.Method+"(\"") {
		return
	}

	parentNameSlice := strings.Split(srcWithourVariableName, "."+router.Method)
	router.ParentName = strings.TrimSpace(parentNameSlice[0])
	return
}

func (router *FindRouters) getGroupPathAndHandlers(src string) {
	// var err error
	// 获取 path And Handler
	groupSlice := strings.Split(src, ".GROUP")
	for k, v := range groupSlice {
		if k == 0 {
			continue
		}

		err := router.getPathAndHandler(v)
		if err != nil {
			continue
		}

	}
	return
}

func (router *FindRouters) getPathAndHandler(paramStr string) (err error) {
	if paramStr[0] != []byte("(")[0] || paramStr[len(paramStr)-1] != []byte(")")[0] {
		return errors.New("format error")
	}

	pathAndHandler := paramStr[1 : len(paramStr)-1]

	phSlice := strings.Split(pathAndHandler, ",")
	if len(phSlice) < 1 {
		return errors.New("format error")
	}

	// path
	pathRaw := strings.TrimSpace(phSlice[0])
	handlerSlice := phSlice[1:]
	if pathRaw[0] != []byte("\"")[0] || pathRaw[len(pathRaw)-1] != []byte("\"")[0] {
		return errors.New("format error")
	}

	router.Path = router.Path + pathRaw[1:len(pathRaw)-1]
	for _, v := range handlerSlice {
		router.HandlersName = append(router.HandlersName, strings.TrimSpace(v))
	}
	return nil
}

// 获取路由 请求方法
func (router *FindRouters) getRouterMethod(src string) (mark bool) {
	if strings.Contains(src, ".GET(") {
		router.Method = "GET"
		return true
	} else if strings.Contains(src, ".POST(") {
		router.Method = "POST"
		return true
	} else if strings.Contains(src, ".PUT(") {
		router.Method = "PUT"
		return true
	} else if strings.Contains(src, ".HEAD(") {
		router.Method = "HEAD"
		return true
	} else if strings.Contains(src, ".DELETE(") {
		router.Method = "DELETE"
		return true
	} else if strings.Contains(src, ".OPTION(") {
		router.Method = "OPTION"
		return true
	} else if strings.Contains(src, ".PATH(") {
		router.Method = "PATH"
		return true
	} else if strings.Contains(src, ".FILE(") {
		router.Method = "FILE"
		return true
	} else if strings.Contains(src, ".PATCH(") {
		router.Method = "PATCH"
		return true
	} else if strings.Contains(src, ".Handle(") {
		router.Method = "Handle"
		return true
	}
	return false
}
