package action

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

// ReadAllFiles 读取所有数据
func ReadAllFiles(rootPath string, serverPath *string, ignoreFolders []string, routers *[]RawRouterStruct, fileSuffix string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ReadAllFiles")
			log.Println(r)
		}
	}()

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if !inSlice(file.Name(), ignoreFolders) {
				ReadAllFiles(rootPath+"/"+file.Name(), serverPath, ignoreFolders, routers, fileSuffix)
			}
		} else {
			if fileSuffix != "" {
				if path.Ext(file.Name()) != fileSuffix {
					continue
				}
			}

			var rawHandlerSlice []string

			body, err := ioutil.ReadFile(rootPath + "/" + file.Name())
			if err != nil {
				return err
			}

			bodySlice := bytes.Split(body, []byte("\n"))

			for _, v := range bodySlice {
				trimBody := string(bytes.TrimSpace(v))
				if len(trimBody) > 4 {
					if trimBody[:2] == "//" {
						if strings.Contains(trimBody, "@ApiServer") {

							// 处理 serverPath
							tmp := strings.Split(trimBody, "@ApiServer")

							if len(tmp) > 1 {
								ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmp[1]))
								if err != nil {
									log.Println("warning: format error —— " + string(v))
									continue
								}

								data := make(map[string]string)
								err = JSON.UnmarshalFromString(ref, &data)
								if err != nil {
									log.Println("warning: format error  —— " + string(v))
									continue
								}
								*serverPath = data["path"]
							}
						} else if strings.Contains(trimBody, "@ApiRouter") {
							// 处理 router
							tmp := strings.Split(trimBody, "@ApiRouter")
							if len(tmp) > 1 {
								ref, err := ParsePMstructToJSONformat(strings.TrimSpace(tmp[1]))
								if err != nil {
									log.Println("warning: format error  ——" + string(v))
									continue
								}
								var router RawRouterStruct
								err = JSON.UnmarshalFromString(ref, &router)
								if err != nil {
									log.Println("warning: format error  ——" + string(v))
									continue
								}

								if inSlice(router.RouterPath, routerPathSlice) {
									log.Println("warning: duplicate Router Path @ApiRouter —— path: \"" + router.RouterPath + "\"")
									continue
								}

								router.Method = strings.ToUpper(router.Method)

								// 生成 router.HandlersName
								handlerSlice := strings.Split(router.RawHandlers, ",")
								for _, v := range handlerSlice {
									router.HandlersName = append(router.HandlersName, strings.TrimSpace(v))
								}

								*routers = append(*routers, router)

								// routerNameSlice = append(routerNameSlice, router.RouterName)
								routerPathSlice = append(routerPathSlice, router.RouterPath)
							}
						} else if strings.Contains(trimBody, "@ApiHandler") || strings.Contains(trimBody, "@ApiQuery") || strings.Contains(trimBody, "@ApiBody") || strings.Contains(trimBody, "@ApiHeader") {
							// 先缓存起来，稍后处理 ApiHandler, ApiQuery, ApiBody, ApiHeader
							rawHandlerSlice = append(rawHandlerSlice, trimBody)
						}
					}
				}
			}
			// 分析 handler, body, query, header
			FormatHandlers(rawHandlerSlice)
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

// ParsePMstructToJSONformat 将(key="id" ...) 转换为 json kv 格式{"key":"id"}
func ParsePMstructToJSONformat(pmStruct string) (string, error) {
	lenPM := len(pmStruct)
	if lenPM < 2 {
		return "", errors.New(pmStruct + "format error")
	}

	if pmStruct[0] != []byte("(")[0] || pmStruct[lenPM-1] != []byte(")")[0] {
		return "", errors.New(pmStruct + "format error")
	}

	var finalStruct []string
	// pmSlice := strings.Split(pmStruct[1:lenPM-1], ",")
	pmSlice := SplitStringsTOKV(pmStruct[1 : lenPM-1])
	for _, v := range pmSlice {
		f, err := formatAPIToJSONKV(v)
		if err != nil {
			log.Println("info: check format ——" + pmStruct)
			continue
		}
		finalStruct = append(finalStruct, f)
	}
	return "{" + strings.Join(finalStruct, ",") + "}", nil
}

// SplitStringsTOKV 将(key="id" ...) 拆分
func SplitStringsTOKV(str string) []string {
	// str = `key="note", value="", desc="备注", type="", src=""`
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var finalSlice []string
	var indexSlice []int
	for k, v := range str {
		if string(v) == "," {
			indexSlice = append(indexSlice, k)
		}
	}

	if len(indexSlice) == 0 {
		finalSlice = append(finalSlice, str)
		return finalSlice
	}

	var finalIndex []int
	lenIndex := len(indexSlice)
	// var mark bool
	for ii := lenIndex - 1; ii >= 0; ii-- {
		if ii > 0 {
			for i := indexSlice[ii] - 1; i >= indexSlice[ii-1]+1; i-- {
				// log.Println(i)
				if string(str[i]) == " " {
					continue
				} else if string(str[i]) == "\"" {
					finalIndex = append(finalIndex, indexSlice[ii])
					break
				} else {
					break
				}
			}
		} else {
			for i := indexSlice[ii] - 1; i >= 0; i-- {
				// log.Println(i)
				if string(str[i]) == " " {
					continue
				} else if string(str[i]) == "\"" {
					finalIndex = append(finalIndex, indexSlice[ii])
					break
				} else {
					break
				}
			}
		}
	}

	lenIndex = len(finalIndex)
	lastIndex := len(str)
	for _, index := range finalIndex {
		if index == finalIndex[lenIndex-1] {
			finalSlice = append(finalSlice, str[index+1:lastIndex])
			finalSlice = append(finalSlice, str[:index])
		} else {
			finalSlice = append(finalSlice, str[index+1:lastIndex])
			lastIndex = index
		}
	}

	return finalSlice
}

// 将api格式转换为 json kv 格式
func formatAPIToJSONKV(KVstr string) (string, error) {
	KVSlice := strings.SplitN(KVstr, "=", 2)
	var jsonStr []string
	if len(KVSlice) < 2 {
		key := strings.ToLower(strings.TrimSpace(KVSlice[0]))
		if key == "" {
			return "", errors.New("format error")
		}

		// 添加 k
		jsonStr = append(jsonStr, "\""+key+"\"")
		// 添加 v
		jsonStr = append(jsonStr, "\"\"")

	} else if len(KVSlice) > 2 {
		return "", errors.New("format error")
	} else {
		// 添加 k
		key := strings.ToLower(strings.TrimSpace(KVSlice[0]))
		if key == "" {
			return "", errors.New("format error")
		}
		jsonStr = append(jsonStr, "\""+key+"\"")

		// 添加 v
		value := strings.TrimSpace(KVSlice[1])
		lenV := len(value)
		if lenV < 2 {
			jsonStr = append(jsonStr, "\"\"")
		} else {
			if value[0] != []byte("\"")[0] || value[lenV-1] != []byte("\"")[0] {
				return "", errors.New("format error")
			}
			jsonStr = append(jsonStr, value)
		}
	}

	return strings.Join(jsonStr, ":"), nil
}
