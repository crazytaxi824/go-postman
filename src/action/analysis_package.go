package action

import (
	"errors"
	"strings"
)

// HandlerPropStruct HandlerPropStruct
type HandlerPropStruct struct {
	FullName        string
	PackageName     string
	VarName         []string
	TypeName        string
	HandlerFuncName string
}

// AnalysisPackage 传入 handlers 里面的 handlers name
// handler = "action.handler1, itemAct.handler2, "
func AnalysisPackage(handlerStruct string) HandlerPropStruct {
	var handler HandlerPropStruct
	handler.FullName = handlerStruct
	// 判断是否是
	packageTypeNameSlice := strings.Split(handlerStruct, ".")
	lenPT := len(packageTypeNameSlice)
	if lenPT < 2 {
		// 没有pacakgename，没有varname，只有控制器name
		handler.HandlerFuncName = packageTypeNameSlice[lenPT-1]
	} else {
		// 判断 package
		b := isPackage(packageTypeNameSlice[0])
		if b {
			// package.var...handlerAct
			handler.PackageName = packageTypeNameSlice[0]
			handler.VarName = packageTypeNameSlice[1 : lenPT-1]
			handler.HandlerFuncName = packageTypeNameSlice[lenPT-1]
		} else {
			// var.var...handlerAct
			handler.VarName = packageTypeNameSlice[0 : lenPT-1]
			handler.HandlerFuncName = packageTypeNameSlice[lenPT-1]
		}
	}
	return handler
}

func isPackage(str string) bool {
	if _, ok := projectFiles[str]; ok {
		return true
	}
	return false
}



// global.ItemsAct.AddHandler
func (hn *HandlerPropStruct) findVarType() error {
	// 

	return nil
}


// findPackageFunction 查找 "package"."Handler" 的handler
func (hn *HandlerPropStruct) findPackageFunction() {
	for _, allFile := range projectFiles[hn.PackageName] {
		var finalContent []string
		var mark bool

		contentSlice := strings.Split(allFile.Content, "\n")
		for _, contentLine := range contentSlice {
			if hn.TypeName == "" {
				if strings.Contains(contentLine, "func") && strings.Contains(contentLine, hn.HandlerFuncName) && !strings.Contains(contentLine, "//") {
					// 添加 @ApiHandler
					key := "// @ApiHandler(name=\"" + hn.FullName + "\")"
					finalContent = append(finalContent, key)
					mark = true
				}
				finalContent = append(finalContent, contentLine)
			} else {
				if strings.Contains(contentLine, "func") && strings.Contains(contentLine, hn.TypeName) && strings.Contains(contentLine, hn.HandlerFuncName) && !strings.Contains(contentLine, "//") {
					// 添加 @ApiHandler
					key := "// @ApiHandler(name=\"" + hn.FullName + "\")"
					finalContent = append(finalContent, key)
					mark = true
				}
				finalContent = append(finalContent, contentLine)
			}
		}
		// 写文件
		if mark {
			fileContent := strings.Join(finalContent, "\n")
			// 写文件
			err := WriteFiles(allFile.FileName, []byte(fileContent))
			if err != nil {
				return
			}
		}
	}
}

// findVarPackage 查找 "var"."Handler" 的handler
// itemsAct.AddHandler
// var itemAct action.ItemAct
func (hn *HandlerPropStruct) findVarPackage(thisFile string) error {
	// 查找 var 的 type struct
	if len(hn.VarName) < 1 {
		return errors.New("format error")
	}

	contentSlice := strings.Split(thisFile, "\n")
	for _, contentLine := range contentSlice {
		if strings.Contains(contentLine, "var "+hn.VarName[0]) {
			typeSlice := strings.Split(contentLine, "var "+hn.VarName[0])
			if len(typeSlice) != 2 {
				return errors.New("format error")
			}

			if strings.Contains(typeSlice[1], "//") {
				typeTmp := strings.Split(typeSlice[1], "//")
				hn.TypeName = strings.TrimSpace(typeTmp[0])
				break
			}
		}
	}

	if hn.TypeName == "" {
		return errors.New("format error")
	}

	if strings.Contains(hn.TypeName, ".") {
		packageType := strings.SplitN(hn.TypeName, ".", 2)
		// TODO
		hn.PackageName = strings.TrimSpace(packageType[0])
		hn.TypeName = strings.TrimSpace(packageType[1])
	}

	hn.findPackageFunction()

	return nil
}