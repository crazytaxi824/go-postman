package action

import (
	"model"
	"strings"
)

// GenFolderAndItemStruct 生成 folder 和 items 给 postman 用
func GenFolderAndItemStruct(groups map[string][]RawRouterStruct) (folders []model.FolderStruct, items []model.RouterStruct) {
	for folderName := range groups {

		if folderName != "undefined" {
			var folder model.FolderStruct
			folder.Name = folderName

			for _, r := range groups[folderName] {
				var router model.RouterStruct
				if strings.TrimSpace(r.RouterName) != "" {
					router.Name = r.RouterName
				} else {
					router.Name = r.RouterPath
				}
				router.Response = make([]interface{}, 0)
				router.Request.URL = r.URL
				router.Request.Method = r.Method
				router.Request.Header = r.Headers
				router.Request.Body = r.Body

				folder.Item = append(folder.Item, router)
			}

			folders = append(folders, folder)
		} else {
			for _, r := range groups[folderName] {
				var router model.RouterStruct
				if strings.TrimSpace(r.RouterName) != "" {
					router.Name = r.RouterName
				} else {
					router.Name = r.RouterPath
				}
				router.Response = make([]interface{}, 0)
				router.Request.URL = r.URL
				router.Request.Method = r.Method
				router.Request.Header = r.Headers
				router.Request.Body = r.Body

				items = append(items, router)
			}
		}
	}

	return
}
