package action

// GroupRouters 将 router 分组
func GroupRouters(routers []RawRouterStruct) map[string][]RawRouterStruct {
	group := make(map[string][]RawRouterStruct)

	for _, r := range routers {
		if r.GroupName == "" {
			group["undefined"] = append(group["undefined"], r)
		} else {
			group[r.GroupName] = append(group[r.GroupName], r)
		}
	}
	return group
}
