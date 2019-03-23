package action

import "strings"

// AnalysisPackage AnalysisPackage
func AnalysisPackage(varName string, fileContent []string) (packageName string) {
	// 检查 import 包
	// if strings.Contains(fileContent, "import \"") || strings.Contains(fileContent, "import (") {
	// 	// 查找 varName 是否存在
	// }
	var importMark bool
	var importStartIndex int
	var importEndIndex int
	for k, contentLine := range fileContent {
		if strings.Contains(contentLine, "import \"") && strings.Contains(contentLine, "\""+varName+"\"") {
			return varName
		}

		if strings.Contains(contentLine, "import (") {
			importMark = true
			importStartIndex = k
		}

		if strings.TrimSpace(contentLine) == ")" {
			importEndIndex = k
		}
	}

	if importMark {
		for i := importStartIndex; i < importEndIndex; i++ {
			if strings.Contains(fileContent[i], "\""+varName+"\"") {
				return varName
			}
		}
	}

	// 检查 var varName 是否存在
	for _, contentLine := range fileContent {
		if strings.Contains(contentLine, "var "+varName) {
			typeName := strings.Split(contentLine, varName)
			return strings.TrimSpace(typeName[1])
		}
	}

	// 获取package name
	return ""
}
