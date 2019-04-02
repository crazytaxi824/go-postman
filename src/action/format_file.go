package action

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"
)

// FindRouters 查找路由用
type FindRouters struct {
	ParentName   string
	VariableName string
	Path         string
	// Handlers     []FindHandlers
	// RootRouter   bool
	HandlersName      []string
	Method            string
	RouterPackageName string
}

// 缓存 文件路径和文件名
var tmpPackageName string

// routerGroups 路由组 缓存
var routerGroups map[string]FindRouters

// 最终的路由，不包含group
var rootRouterGroups []FindRouters

// ProjectFiles 缓存file，key-package name;
var ProjectFiles map[string][]AllFiles

// AllFiles 文件内容
type AllFiles struct {
	FileName   string
	Content    string
	FormatMark bool
}

// ReformatFile 逐行遍历，添加 Api 文件
func ReformatFile(rootPath string, ignoreFolders []string, fileSuffix string) error {
	routerGroups = make(map[string]FindRouters)

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	// var fileName string
	for _, file := range files {
		if file.IsDir() {
			if !inSlice(file.Name(), ignoreFolders) {
				ReformatFile(rootPath+"/"+file.Name(), ignoreFolders, fileSuffix)
			}
		} else {
			if fileSuffix != "" {
				if path.Ext(file.Name()) != fileSuffix {
					continue
				}
			}

			tmpPackageName = ""

			// 最终要写入文件的内容
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
				// 获取所有文件内容存入 projectFile 中
				if strings.Contains(str, "package") && !strings.Contains(str, "//") {
					// 获取 pacakge name
					packageNameSlice := strings.Split(str, "package")
					// 缓存文件 package 名称
					tmpPackageName = strings.TrimSpace(packageNameSlice[1])
				}

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
							// fileName = file.Name()
							mark = true
						} else if strings.TrimSpace(bodySlice[k-1]) != strings.TrimSpace(apiStr) {
							finalFile = append(finalFile, apiStr)
							// fileName = file.Name()
							mark = true
						}
					} else {
						finalFile = append(finalFile, apiStr)
						// fileName = file.Name()
						mark = true
					}
				}

				finalFile = append(finalFile, str)
			}

			// 存入 fileName 和 全部数据
			if tmpPackageName != "" {
				var fileContent AllFiles
				fileContent.FileName = filePath
				fileContent.Content = strings.Join(finalFile, "\n")
				fileContent.FormatMark = mark

				ProjectFiles[tmpPackageName] = append(ProjectFiles[tmpPackageName], fileContent)
			} else {
				// log.Println("packageName does not exist: " + file.Name())
				continue
			}

		}
	}
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
	} else if strings.Contains(src, ".Header.Get(\"") && strings.Contains(src, "\")") && !strings.Contains(src, "//") {
		i := strings.Index(src, ".Header.Get(\"")
		f := strings.Index(src, "\")")
		key := strings.TrimSpace(src[i+13 : f])

		// @ApiHeader(key="header",value="value",desc="header description")
		header := "// @ApiHeader(key=\"" + key + "\", desc=\"\", value=\"\")"
		return header, nil

		// } else if strings.Contains(src, ".GROUP(\"") && !strings.Contains(src, "//") {
	} else if !strings.Contains(src, "//") {
		r, err := parseRouterGroupProperties(src)
		if err != nil {
			return "", err
		}

		if r != "" {
			return r, nil
		}
	}
	return "", nil
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

// eg: itemAct.GROUP("/get", action.handler).GROUP("/get2", action.handler2).GET("/get3", action.handler3)
func (router *FindRouters) genRouterAPI(src string) (string, error) {
	// 缓存路由 文件名
	router.RouterPackageName = tmpPackageName
	// .GET 分割
	routerSlice := strings.Split(src, "."+router.Method)
	// itemAct.GROUP("/get", action.handler).GROUP("/get2", action.handler2)
	groupsRaw := routerSlice[0]
	// ("/get3", action.handler3)
	routerRaw := routerSlice[1]

	// 分析 ("/get3", action.handler3)
	err := router.getPathAndHandler(strings.TrimSpace(routerRaw))
	if err != nil {
		return "", err
	}

	// 如果不包含 .GROUP
	if !strings.Contains(groupsRaw, ".GROUP") {
		router.getRouterParentName(src)
		// 向上查找 parent
		router.findingParent(router.ParentName)

		// 添加到 rootRouterGroups
		rootRouterGroups = append(rootRouterGroups, *router)

		apiStr := "// @ApiRouter(path=\"" + router.Path + "\", desc=\"\", method=\"" + router.Method + "\", group=\"" + router.ParentName + "\", handlers=\"" + strings.Join(router.HandlersName, ",") + "\")"

		return apiStr, nil
	}

	// itemAct.GROUP("/get", action.handler).GROUP("/get2", action.handler2)
	var group FindRouters

	// [itemAct ("/get", action.handler) ("/get2", action.handler2)]
	groupsSlice := strings.Split(groupsRaw, ".GROUP")
	for k, v := range groupsSlice {
		// 第一部分不处理
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

	apiStr := "// @ApiRouter(path=\"" + router.Path + "\", desc=\"\", method=\"" + router.Method + "\", group=\"" + group.ParentName + "\", handlers=\"" + strings.Join(router.HandlersName, ",") + "\")"

	return apiStr, nil
}

// 递归查找上级 group
func (router *FindRouters) findingParent(parentName string) {
	router.Path = routerGroups[parentName].Path + router.Path
	router.HandlersName = append(router.HandlersName, routerGroups[parentName].HandlersName...)
	if routerGroups[parentName].ParentName != "" {
		router.findingParent(routerGroups[parentName].ParentName)
	}
}

// eg: itemAct.GROUP("/get", action.handler).GROUP("/get2", action.handler2).GROUP("/get3", action.handler3)
func (router *FindRouters) getGroupParentName(srcWithourVariableName string) {
	// eg: d.Router.GROUP("/get", action.handler)
	if strings.Contains(srcWithourVariableName, ".Router.GROUP(\"") {
		return
	}

	// 获取到 itemAct
	parentNameSlice := strings.Split(srcWithourVariableName, ".GROUP")
	router.ParentName = strings.TrimSpace(parentNameSlice[0])
	return
}

// eg: itemAct.GET("/get", action.handler)
func (router *FindRouters) getRouterParentName(srcWithourVariableName string) {
	// eg: d.Router.GET("/get", action.handler)
	if strings.Contains(srcWithourVariableName, ".Router."+router.Method+"(\"") {
		return
	}

	// 获取到 itemAct
	parentNameSlice := strings.Split(srcWithourVariableName, "."+router.Method)
	router.ParentName = strings.TrimSpace(parentNameSlice[0])
	return
}

// eg: itemAct.GROUP("/get", action.handler).GROUP("/get2", action.handler2).GROUP("/get3", action.handler3)
func (router *FindRouters) getGroupPathAndHandlers(src string) {
	// var err error
	// 获取 path And Handler
	groupSlice := strings.Split(src, ".GROUP")
	for k, v := range groupSlice {
		// 跳过第一部分
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

// 获取单个路由的路径和处理器名称，eg: ("/get", action.handler)
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
