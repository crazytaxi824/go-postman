package model

// PostmanStruct PostmanStruct
type PostmanStruct struct {
	Info InfoStruct    `json:"info"`
	Item []interface{} `json:"item"` //item可以为 FolderStruct, 也可以为 SingleRequest
}

// InfoStruct InfoStruct
type InfoStruct struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

// FolderStruct FolderStruct
type FolderStruct struct {
	Name               string        `json:"name"`
	Item               []interface{} `json:"item"`                           // item可以为 FolderStruct, 也可以为 SingleRequest
	PostmanIsSubFolder bool          `json:"_postman_isSubFolder,omitempty"` // 注意json字段不要写错了 _postman_isSubFolder
}

// RouterStruct RouterStruct
type RouterStruct struct {
	Name     string        `json:"name"`
	Request  RequestStruct `json:"request"`
	Response []interface{} `json:"response"`
}

// RequestStruct RequestStruct
type RequestStruct struct {
	Method string         `json:"method"`
	Header []HeaderStruct `json:"header"`
	Body   interface{}    `json:"body"`
	URL    URLStruct      `json:"url"`
}

// URLStruct URLStruct
type URLStruct struct {
	Raw      string        `json:"raw"`
	Protocol string        `json:"protocol"`
	Host     []string      `json:"host"`
	Port     string        `json:"port,omitempty"`
	Path     []string      `json:"path"`
	Query    []QueryStruct `json:"query,omitempty"`
}

// QueryStruct QueryStruct
type QueryStruct struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// ModeDataFileStruct ModeDataFileStruct
type ModeDataFileStruct struct {
	Key         string `json:"key"`
	Description string `json:"description,omitempty"` // 描述
	Type        string `json:"type"`                  // file
	Src         string `json:"src"`                   // file 地址
}

// ModeDataTextStruct ModeDataTextStruct
type ModeDataTextStruct struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"` // 描述
	Type        string `json:"type"`                  // text
}

// HeaderStruct HeaderStruct
type HeaderStruct struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}
