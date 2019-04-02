package action

import (
	"model"
)

// HandlersToRouters 将 handler 的所有参数传到对应的 router 中
func (router *RawRouterStruct) HandlersToRouters() {

	for _, handlerName := range router.HandlersName {
		for _, rawHandler := range HandlerMap[handlerName] {

			// TODO 判断是否有重复的key

			router.Headers = append(router.Headers, rawHandler.Headers...)
			if router.Method == "GET" {

				router.Querys = append(router.Querys, rawHandler.Querys...)
				for _, v := range rawHandler.Texts {
					var query model.QueryStruct
					query.Key = v.Key
					query.Description = v.Description
					query.Value = v.Value

					router.Querys = append(router.Querys, query)
				}

			} else {

				router.Files = append(router.Files, rawHandler.Files...)
				router.Texts = append(router.Texts, rawHandler.Texts...)
				router.Querys = append(router.Querys, rawHandler.Querys...)
			}
		}
	}

}
