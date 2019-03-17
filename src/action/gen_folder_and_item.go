package action

import "model"

// GenFolderAndItemStruct 生成folder和items
func GenFolderAndItemStruct(groups map[string][]RawRouterStruct) (folders []model.FolderStruct, items []model.RouterStruct) {
	for folderName := range groups {

		if folderName != "undefined" {
			var folder model.FolderStruct
			folder.Name = folderName

			for _, r := range groups[folderName] {
				var router model.RouterStruct
				router.Name = r.RouterName
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
				router.Name = r.RouterName
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
