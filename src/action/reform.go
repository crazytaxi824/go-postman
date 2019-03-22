package action

import (
	"bytes"
	"io/ioutil"
)

// ReformFile 逐行遍历，添加 Api 文件
func ReformFile(rootPath string, ignoreFolders []string) error {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if !inSlice(file.Name(), ignoreFolders) {
				ReformFile(rootPath+"/"+file.Name(), ignoreFolders)
			}
		} else {
			var finalFile [][]byte
			// 是否需要写文件
			mark := false

			filePath := rootPath + "/" + file.Name()
			body, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			bodySlice := bytes.Split(body, []byte("\n"))

			for k, b := range bodySlice {
				if k != 0 {
					if bytes.Contains(b, []byte(".QueryValue(\"")) && bytes.Contains(b, []byte("\")")) && !bytes.Contains(bodySlice[k-1], []byte("@ApiQuery(")) && !bytes.Contains(b, []byte("//")) {

						i := bytes.Index(b, []byte(".QueryValue(\""))
						f := bytes.Index(b, []byte("\")"))
						key := []byte(b[i+13 : f])
						var queryBytes [][]byte
						queryBytes = append(queryBytes, []byte("// @ApiQuery(key=\""))
						queryBytes = append(queryBytes, bytes.TrimSpace([]byte(key)))
						queryBytes = append(queryBytes, []byte("\",desc= \"\", value=\"\")"))
						query := bytes.Join(queryBytes, nil)

						finalFile = append(finalFile, query)

						mark = true

					} else if bytes.Contains(b, []byte(".FormValue(\"")) && bytes.Contains(b, []byte("\")")) && !bytes.Contains(bodySlice[k-1], []byte("@ApiBody(")) && !bytes.Contains(b, []byte("//")) {
						i := bytes.Index(b, []byte(".FormValue(\""))
						f := bytes.Index(b, []byte("\")"))

						key := []byte(b[i+12 : f])
						var bodyBytes [][]byte
						bodyBytes = append(bodyBytes, []byte("// @ApiBody(key=\""))
						bodyBytes = append(bodyBytes, bytes.TrimSpace([]byte(key)))
						bodyBytes = append(bodyBytes, []byte("\",desc=\"\", value=\"\")"))
						body := bytes.Join(bodyBytes, nil)

						finalFile = append(finalFile, body)
						mark = true
					}
				} else {
					if bytes.Contains(b, []byte(".QueryValue(\"")) && bytes.Contains(b, []byte("\")")) && !bytes.Contains(b, []byte("//")) {
						i := bytes.Index(b, []byte(".QueryValue(\""))
						f := bytes.Index(b, []byte("\")"))
						key := []byte(b[i+13 : f])
						var queryBytes [][]byte
						queryBytes = append(queryBytes, []byte("// @ApiQuery(key=\""))
						queryBytes = append(queryBytes, bytes.TrimSpace([]byte(key)))
						queryBytes = append(queryBytes, []byte("\",desc= \"\", value=\"\")"))
						query := bytes.Join(queryBytes, nil)

						finalFile = append(finalFile, query)

						mark = true
					} else if bytes.Contains(b, []byte(".FormValue(\"")) && bytes.Contains(b, []byte("\")")) && !bytes.Contains(b, []byte("//")) {
						i := bytes.Index(b, []byte(".FormValue(\""))
						f := bytes.Index(b, []byte("\")"))

						key := []byte(b[i+12 : f])
						var bodyBytes [][]byte
						bodyBytes = append(bodyBytes, []byte("// @ApiBody(key=\""))
						bodyBytes = append(bodyBytes, bytes.TrimSpace([]byte(key)))
						bodyBytes = append(bodyBytes, []byte("\",desc=\"\", value=\"\")"))
						body := bytes.Join(bodyBytes, nil)

						finalFile = append(finalFile, body)
						mark = true
					}
				}
				finalFile = append(finalFile, b)
			}

			// 写文件
			if mark {
				fileContent := bytes.Join(finalFile, []byte("\n"))
				err = WriteFiles(filePath, fileContent)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
