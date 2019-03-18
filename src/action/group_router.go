package action

import (
	"strings"
)

// GroupRouters 将 router 按照 gourp name 分组，否则不会被放到文件夹中
func GroupRouters(routers []RawRouterStruct) map[string][]RawRouterStruct {
	group := make(map[string][]RawRouterStruct)

	for _, r := range routers {
		if strings.TrimSpace(r.GroupName) == "" {
			group["undefined"] = append(group["undefined"], r)
		} else {
			group[r.GroupName] = append(group[r.GroupName], r)
		}
	}
	return group
}
