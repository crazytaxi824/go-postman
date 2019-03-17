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
								ref, err := ParsePMstructToJSONformat(tmp[1])
								if err != nil {
									log.Println("warning: 格式错误 —— " + string(v))
									continue
								}

								data := make(map[string]string)
								err = JSON.UnmarshalFromString(ref, &data)
								if err != nil {
									log.Println("warning: 格式错误 —— " + string(v))
									continue
								}
								*serverPath = data["path"]
							}
						} else if strings.Contains(trimBody, "@pmRouter") {
							// 处理 router
							tmp := strings.Split(trimBody, "@pmRouter")
							if len(tmp) > 1 {
								ref, err := ParsePMstructToJSONformat(tmp[1])
								if err != nil {
									log.Println("warning: 格式错误 ——" + string(v))
									continue
								}
								var router RawRouterStruct
								err = JSON.UnmarshalFromString(ref, &router)
								if err != nil {
									log.Println("warning: 格式错误 ——" + string(v))
									continue
								}

								// 判断router 是否存在，是否重名
								if inSlice(router.RouterName, routerNameSlice) {
									log.Println("warning: 项目中有两个相同名字的路由 @pmRouter —— " + router.RouterName)
									continue
								}

								router.Method = strings.ToUpper(router.Method)

								*routers = append(*routers, router)
								routerNameSlice = append(routerNameSlice, router.RouterName)
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

// ParsePMstructToJSONformat ParsePMstructToJSONformat
func ParsePMstructToJSONformat(pmStruct string) (string, error) {
	lenPM := len(pmStruct)
	if lenPM < 2 {
		return "", errors.New(pmStruct + "格式错误")
	}

	if pmStruct[0] != []byte("(")[0] || pmStruct[lenPM-1] != []byte(")")[0] {
		return "", errors.New(pmStruct + "格式错误")
	}

	var finalStruct []string
	// pmSlice := strings.Split(pmStruct[1:lenPM-1], ",")
	pmSlice := SplitStringsTOKV(pmStruct[1 : lenPM-1])
	for _, v := range pmSlice {
		f, err := parseKV(v)
		if err != nil {
			log.Println("info: 请检查格式 ——" + pmStruct)
			continue
		}
		finalStruct = append(finalStruct, f)
	}
	return "{" + strings.Join(finalStruct, ",") + "}", nil
}

// SplitStringsTOKV SplitStringsTOKV
func SplitStringsTOKV(str string) []string {
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

func parseKV(KVstr string) (string, error) {
	KVSlice := strings.SplitN(KVstr, "=", 2)
	var jsonStr []string
	if len(KVSlice) < 2 {
		key := strings.ToLower(strings.TrimSpace(KVSlice[0]))
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
		key := strings.ToLower(strings.TrimSpace(KVSlice[0]))
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
