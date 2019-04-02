package action

import "log"

// HandlersToRouters 将 handler 的所有参数传到对应的 router 中
func (router *RawRouterStruct) HandlersToRouters() {
	for _, handlerName := range router.HandlersName {
		log.Println(handlerName)
		for _, rawHandler := range HandlerMap[handlerName] {
			// log.Println(rawHandler)

			// TODO 判断是否有重复的key
			router.Files = append(router.Files, rawHandler.Files...)
			router.Texts = append(router.Texts, rawHandler.Texts...)

			router.Headers = append(router.Headers, rawHandler.Headers...)
			router.Querys = append(router.Querys, rawHandler.Querys...)
		}
	}
}
