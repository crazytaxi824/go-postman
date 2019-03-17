package action

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// ReadAllFiles 读取所有数据
func ReadAllFiles(rootPath string, serverPath *string, ignoreFolders []string, routers *[]RawRouterStruct, rawHandlerSlice *[]string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		if file.IsDir() {
			if !inSlice(file.Name(), ignoreFolders) {
				ReadAllFiles(rootPath+"/"+file.Name(), serverPath, ignoreFolders, routers, rawHandlerSlice)
			}
		} else {
			body, err := ioutil.ReadFile(rootPath + "/" + file.Name())
			if err != nil {
				return err
			}

			bodySlice := bytes.Split(body, []byte("\n"))

			for _, v := range bodySlice {
				trimBody := string(bytes.TrimSpace(v))
				if len(trimBody) > 4 {
					if trimBody[:2] == "//" {
						if strings.Contains(trimBody, "@pmServer") {

							// 处理 serverPath
							tmp := strings.Split(trimBody, "@pmServer")

							if len(tmp) > 1 {
								ref, err := ParsePMstruct(tmp[1])
								if err != nil {
									continue
								}

								data := make(map[string]string)
								// err = JSON.UnmarshalFromString(tmp[1], &data)
								err = JSON.UnmarshalFromString(ref, &data)
								if err != nil {
									// log.Println(string(v) + "pmServer 格式错误")
									continue
								}
								// if _, ok := data["path"]; !ok {
								// 	log.Println("pmServer 没有\"path\"参数")
								// }
								*serverPath = data["path"]
							}
						} else if strings.Contains(trimBody, "@pmRouter") {

							// 处理 router
							tmp := strings.Split(trimBody, "@pmRouter")
							if len(tmp) > 1 {
								ref, err := ParsePMstruct(tmp[1])
								if err != nil {
									continue
								}
								var router RawRouterStruct
								// err = JSON.UnmarshalFromString(tmp[1], &router)
								err = JSON.UnmarshalFromString(ref, &router)

								router.Method = strings.ToUpper(router.Method)
								if err != nil {
									// log.Println(string(v) + " —— 格式错误")
									continue
								}
								*routers = append(*routers, router)
							}
						} else if strings.Contains(trimBody, "@pmHandler") || strings.Contains(trimBody, "@pmQuery") || strings.Contains(trimBody, "@pmBody") || strings.Contains(trimBody, "@pmHeader") {
							// 处理 pmHandler, pmQuery, pmBody, pmHeader
							*rawHandlerSlice = append(*rawHandlerSlice, trimBody)
						}
					}
				}
			}
		}
	}
	return nil
}

// WriteFiles 写文件
func WriteFiles(filename string, content []byte) error {
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func inSlice(s string, ss []string) bool {
	for k := range ss {
		if ss[k] == s {
			return true
		}
	}
	return false
}

// ParsePMstruct ParsePMstruct
func ParsePMstruct(pmStruct string) (string, error) {
	// @ApiResponse(code = CommonStatus.EXCEPTION, message = "服务器内部异常")
	// @pmRouter(name="添加文章", method="Post", path="/m/article/add", group="文章")
	lenPM := len(pmStruct)
	if lenPM < 2 {
		return "", errors.New(pmStruct + "格式错误")
	}

	if pmStruct[0] != []byte("(")[0] || pmStruct[lenPM-1] != []byte(")")[0] {
		return "", errors.New(pmStruct + "格式错误")
	}

	var finalStruct []string
	pmSlice := strings.Split(pmStruct[1:lenPM-1], ",")
	for _, v := range pmSlice {
		f, err := parseKV(v)
		if err != nil {
			continue
		}
		finalStruct = append(finalStruct, f)
	}
	return "{" + strings.Join(finalStruct, ",") + "}", nil
}

func parseKV(KVstr string) (string, error) {
	KVSlice := strings.SplitN(KVstr, "=", 2)
	var jsonStr []string
	if len(KVSlice) < 2 {
		key := strings.TrimSpace(KVSlice[0])
		if key == "" {
			return "", errors.New("格式错误")
		}

		// 添加 k
		jsonStr = append(jsonStr, "\""+key+"\"")
		// 添加 v
		jsonStr = append(jsonStr, "\"\"")

	} else if len(KVSlice) > 2 {
		return "", errors.New("格式错误")
	} else {
		// 添加 k
		key := strings.TrimSpace(KVSlice[0])
		if key == "" {
			return "", errors.New("格式错误")
		}
		jsonStr = append(jsonStr, "\""+key+"\"")

		// 添加 v
		value := strings.TrimSpace(KVSlice[1])
		lenV := len(value)
		if lenV < 2 {
			jsonStr = append(jsonStr, "\"\"")
		} else {
			if value[0] != []byte("\"")[0] || value[lenV-1] != []byte("\"")[0] {
				return "", errors.New("格式错误")
			}
			jsonStr = append(jsonStr, value)
		}
	}

	return strings.Join(jsonStr, ":"), nil
}
