## go-postman

分析源代码中的注释，解析成postman的json文件，方便导入postman。
swagger风格。

### 事例

服务路径
```go
func server(){
  // @ApiServer(path = "http://127.0.0.1:18080")
  srv.ListenAndServe := ...
}
```
路由
```go
func Router(){
  ...
  // @ApiRouter(name="添加文章", method="Post", path="/article/add", group="文章")
  article.Get("/add", articleAct.AddArticle)
  
  // @ApiRouter(name= "编辑文章",method= "Post",path= "/m/article/edit",group= "文章")
  article.Get("/edit", articleAct.EditArticle)
  
  // @ApiRouter(name= "文章列表",method= "Get",path= "/m/article/list",group= "文章")  
  article.Get("/list", articleAct.ListArticle)
  ...
}
```
控制器
```go
// 处理器名称必须和路由对应，否则会被抛弃
// @ApiHandler(name= "文章列表")
func ListHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiQuery(key= "classID", value="123", desc= "文章类型id")
  cID := req.URL.Query().Get("classID")
  ...
}

// 处理器名称必须和路由对应，否则会被抛弃
// @ApiHandler(name= "编辑文章")
func EditHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiQuery(key= "id", value="123", desc= "文章id")
  ArticleID := req.URL.Query().Get("id")
  
  // @ApiBody(key="title",desc="文章标题")
  articleTitle := req.PostFormValue("title")
  
  // @ApiBody(key="content",desc="文章内容")
  articleContent := req.PostFormValue("content")
  
  // @ApiBody(key="author", value="xxx", desc="作者")
  articleAuthor := req.PostFormValue("author")
  
  // @ApiBody(key="author",desc="图片",type="file",src="/eee.png")
  ...
}

// 处理器名称必须和路由对应，否则会被抛弃
// @ApiHandler(name= "添加文章")
func AddHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiHeader(key="Content-Type",value="application/x-www-form-urlencoded",desc="header描述")
  // @ApiHeader(key="Content-Type",value="application/json",desc="header描述")
  
  // @ApiQuery(key= "classID", value="123", desc= "文章类型id")
  
  // @ApiBody(key="title",desc="文章标题")
  // @ApiBody(key="content",desc="文章内容")
  // @ApiBody(key="author",desc="作者")
  ...
}

```
-----

### 命令行工具
```bash
$ gpm
```
参数
```bash
  -i string
    	不读取指定文件夹名称下的所有文件, 用 | 分隔多个文件夹 (default "vendor")
  -o string
    	输出json文件的路径和名称 (default "./newPostman.json")
  -p string
    	指定项目路径，默认从src文件夹下开始读取 (default "./src")
```

