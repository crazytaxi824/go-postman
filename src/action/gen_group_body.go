package action

// GenGroupBody 生成body
func GenGroupBody(groupsPointer *map[string][]RawRouterStruct) {
	groups := *groupsPointer
	for k := range groups {
		for index := range groups[k] {
			finalBodyData := make(map[string]interface{})

			lenFile := len(groups[k][index].Files)

			if lenFile+len(groups[k][index].Texts) == 0 {
				finalBodyData["mode"] = "raw"
				finalBodyData["raw"] = ""

				groups[k][index].Body = finalBodyData
				continue
			}

			var bodyData []interface{}
			mark := false
			if lenFile > 0 {
				mark = true
			}

			for _, v := range groups[k][index].Files {
				bodyData = append(bodyData, v)
			}

			for _, v := range groups[k][index].Texts {
				bodyData = append(bodyData, v)
			}

			if mark {
				finalBodyData["mode"] = "formdata"
				finalBodyData["formdata"] = bodyData
			} else {
				finalBodyData["mode"] = "urlencoded"
				finalBodyData["urlencoded"] = bodyData
			}

			groups[k][index].Body = finalBodyData
		}
	}
}
