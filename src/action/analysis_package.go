package action

import (
	"errors"
	"log"
	"strings"
)

// HandlerPropStruct HandlerPropStruct
type HandlerPropStruct struct {
	FullName           string
	VarStructs         []VarStruct
	HandlerFuncName    string
	HandlerPackageName string
}

// VarStruct 变量类型
type VarStruct struct {
	VarPackageName string
	VarName        string
	VarType        TypeStruct
}

// TypeStruct 类型
type TypeStruct struct {
	typePackageName string
	typeName        string
	// fileName        string
}

// 缓存handler, 去重
var projectHandler map[string]HandlerPropStruct

// AnalysisFindRouter 分析rootRouterGroups，缓存 projectHandler
func AnalysisFindRouter() error {
	projectHandler = make(map[string]HandlerPropStruct)
	for _, router := range rootRouterGroups {
		for _, handlerName := range router.HandlersName {
			if _, ok := projectHandler[handlerName]; ok {
				continue
			}
			projectHandler[handlerName] = AnalysisRawPackage(handlerName, router.RouterPackageName)
		}
	}

	for _, handler := range projectHandler {
		err := handler.analysisHandlerVarStruct()
		if err != nil {
			return err
		}
	}

	// 写文件

	for _, handler := range projectHandler {
		if len(handler.VarStructs) < 1 {
			handler.analysisPackage()
		} else {
			handler.analysisLastVar()
		}
	}

	return nil
}

// AnalysisRawPackage 传入 handlers 里面的 handlers name
// handler = "action.handler1, itemAct.handler2, "
func AnalysisRawPackage(handlerName, packageName string) HandlerPropStruct {
	var handler HandlerPropStruct
	handler.FullName = handlerName
	handler.HandlerPackageName = packageName
	// 判断是否是
	packageTypeNameSlice := strings.Split(handlerName, ".")
	lenPT := len(packageTypeNameSlice)
	if lenPT < 2 {
		// 没有pacakgename，没有varname，只有控制器name
		handler.HandlerFuncName = strings.TrimSpace(packageTypeNameSlice[lenPT-1])

	} else {
		// 判断 package
		if _, ok := ProjectFiles[packageTypeNameSlice[0]]; ok {
			// package 存在
			if lenPT < 3 {
				// package.handlerAct
				handler.HandlerPackageName = strings.TrimSpace(packageTypeNameSlice[0])

			} else {
				// package.Var...handlerAct
				for i := 1; i < lenPT-1; i++ {
					if i == 1 {
						var varStruct VarStruct
						varStruct.VarPackageName = strings.TrimSpace(packageTypeNameSlice[0])
						varStruct.VarName = strings.TrimSpace(packageTypeNameSlice[1])
						handler.VarStructs = append(handler.VarStructs, varStruct)
					} else {
						var varStruct VarStruct
						varStruct.VarName = strings.TrimSpace(packageTypeNameSlice[i])
						handler.VarStructs = append(handler.VarStructs, varStruct)
					}
				}
			}
			handler.HandlerFuncName = strings.TrimSpace(packageTypeNameSlice[lenPT-1])
			// handler.VarPackageName = packageTypeNameSlice[0]
			// handler.VarName = packageTypeNameSlice[1 : lenPT-1]
		} else {
			// package 不存在
			// var.var...handlerAct
			for i := 0; i < lenPT-1; i++ {
				if i == 0 {
					var varStruct VarStruct
					varStruct.VarPackageName = strings.TrimSpace(handler.HandlerPackageName)
					varStruct.VarName = strings.TrimSpace(packageTypeNameSlice[0])
					handler.VarStructs = append(handler.VarStructs, varStruct)
				} else {
					var varStruct VarStruct
					varStruct.VarName = strings.TrimSpace(packageTypeNameSlice[i])
					handler.VarStructs = append(handler.VarStructs, varStruct)
				}
			}
			handler.HandlerFuncName = strings.TrimSpace(packageTypeNameSlice[lenPT-1])
		}
	}
	// log.Println(handler.FullName, handler.HandlerPackageName, handler.HandlerFuncName, handler.VarStructs)
	return handler
}

func (vs *VarStruct) findFirstVarTypeStruct() error {
	var files []AllFiles
	files = ProjectFiles[vs.VarPackageName]

	if len(files) < 1 {
		return errors.New("cannot find package files")
	}

	for _, file := range files {
		bodySlice := strings.Split(file.Content, "\n")
		for _, bodyLine := range bodySlice {
			if strings.Contains(bodyLine, "var "+vs.VarName) && !strings.Contains(bodyLine, "//") {
				tmpTypeName := strings.Split(bodyLine, "var "+vs.VarName)

				varTypeName := strings.TrimSpace(tmpTypeName[1])
				if strings.Contains(varTypeName, ".") {
					tmp := strings.Split(varTypeName, ".")

					vs.VarType.typePackageName = strings.TrimSpace(tmp[0])
					vs.VarType.typeName = strings.TrimSpace(tmp[1])
					return nil
				}

				vs.VarType.typePackageName = strings.TrimSpace(vs.VarPackageName)
				vs.VarType.typeName = strings.TrimSpace(varTypeName)
				return nil
			}
		}

	}

	return errors.New("cannot find var Type")
}

func (vs *VarStruct) findRestVarTypeStruct(ts TypeStruct) error {
	var files []AllFiles
	files = ProjectFiles[ts.typePackageName]
	if len(files) < 1 {
		return errors.New("cannot find package files")
	}

	var mark bool
	for _, file := range files {
		bodySlice := strings.Split(file.Content, "\n")
		for _, bodyLine := range bodySlice {
			if !mark {
				if strings.Contains(bodyLine, "type "+ts.typeName+" struct {") && !strings.Contains(bodyLine, "//") {
					mark = true
				}
			} else {
				if strings.Contains(bodyLine, vs.VarName) && !strings.Contains(bodyLine, "//") {
					tmp := strings.Split(bodyLine, vs.VarName)
					typeName := strings.TrimSpace(tmp[1])
					if strings.Contains(typeName, ".") {
						tmpTypeStruct := strings.Split(typeName, ".")
						vs.VarType.typePackageName = strings.TrimSpace(tmpTypeStruct[0])
						vs.VarType.typeName = strings.TrimSpace(tmpTypeStruct[1])
						return nil
					}

					vs.VarType.typePackageName = ts.typePackageName
					vs.VarType.typeName = strings.TrimSpace(typeName)
					return nil
				}
			}
		}
	}

	return errors.New("cannot find var Type")
}

func (h *HandlerPropStruct) analysisHandlerVarStruct() error {
	for k := range h.VarStructs {
		if k == 0 {
			err := h.VarStructs[k].findFirstVarTypeStruct()
			if err != nil {
				log.Println(err.Error())
				return err
			}
		} else {
			err := h.VarStructs[k].findRestVarTypeStruct(h.VarStructs[k-1].VarType)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
	}
	return nil
}

func (h *HandlerPropStruct) analysisLastVar() {
	lenType := len(h.VarStructs)

	lastPackageName := h.VarStructs[lenType-1].VarType.typePackageName
	lastTypeName := h.VarStructs[lenType-1].VarType.typeName

	for kk := range ProjectFiles[lastPackageName] {
		bodySlice := strings.Split(ProjectFiles[lastPackageName][kk].Content, "\n")
		for k, bodyLine := range bodySlice {
			if k == 0 {
				continue
			}

			if strings.Contains(bodyLine, "func ") && strings.Contains(bodyLine, lastTypeName+") "+h.HandlerFuncName+"(") && strings.Contains(bodyLine, "httpdispatcher.Context) error {") && !strings.Contains(bodyLine, "//") {

				key := "// @ApiHandler(name=\"" + h.FullName + "\")"
				if strings.Contains(bodySlice[k-1], "//") && strings.Contains(bodySlice[k-1], "@ApiHandler") {
					if !strings.Contains(bodySlice[k-1], key) {
						tmpBodyFront := bodySlice[: k-1 : k-1]
						tmpBodyEnd := bodySlice[k:]
						final := tmpBodyFront
						final = append(final, key)
						final = append(final, tmpBodyEnd...)
						ProjectFiles[lastPackageName][kk].Content = strings.Join(final, "\n")
						ProjectFiles[lastPackageName][kk].FormatMark = true
						return
					}
					return
				}

				tmpBodyFront := bodySlice[:k:k]
				tmpBodyEnd := bodySlice[k:]
				final := tmpBodyFront
				final = append(final, key)
				final = append(final, tmpBodyEnd...)

				ProjectFiles[lastPackageName][kk].Content = strings.Join(final, "\n")
				ProjectFiles[lastPackageName][kk].FormatMark = true
				return
			}
		}
	}
}

func (h *HandlerPropStruct) analysisPackage() {
	for kk := range ProjectFiles[h.HandlerPackageName] {

		bodySlice := strings.Split(ProjectFiles[h.HandlerPackageName][kk].Content, "\n")
		for k, bodyLine := range bodySlice {
			if k == 0 {
				continue
			}

			if strings.Contains(bodyLine, "func "+h.HandlerFuncName+"(") && strings.Contains(bodyLine, "httpdispatcher.Context) error {") && !strings.Contains(bodyLine, "//") {

				key := "// @ApiHandler(name=\"" + h.FullName + "\")"
				if strings.Contains(bodySlice[k-1], "//") && strings.Contains(bodySlice[k-1], "@ApiHandler") {
					if !strings.Contains(bodySlice[k-1], key) {
						tmpBodyFront := bodySlice[: k-1 : k-1]
						tmpBodyEnd := bodySlice[k:]
						final := tmpBodyFront
						final = append(final, key)
						final = append(final, tmpBodyEnd...)
						ProjectFiles[h.HandlerPackageName][kk].Content = strings.Join(final, "\n")
						ProjectFiles[h.HandlerPackageName][kk].FormatMark = true
						return
					}
					return
				}

				tmpBodyFront := bodySlice[:k:k]
				tmpBodyEnd := bodySlice[k:]
				final := tmpBodyFront
				final = append(final, key)
				final = append(final, tmpBodyEnd...)
				ProjectFiles[h.HandlerPackageName][kk].Content = strings.Join(final, "\n")
				ProjectFiles[h.HandlerPackageName][kk].FormatMark = true
				return

			}
		}
	}
}
