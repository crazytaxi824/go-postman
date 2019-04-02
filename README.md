## go-postman

分析源代码中的注释，解析成postman的json文件，方便导入postman。swagger风格。

Parse comments from source code, and reformat into PostMan json file. swagger style.

### 事例 example

#### 服务路径 server path

```go
func server(){
  // @ApiServer(path = "http://127.0.0.1:18080")
  srv.ListenAndServe := ...
}
```

#### 路由 router

路由中的 handler 可以添加多个处理器，用","分隔

"handlers" in @ApiRouter could be mutli handlers' names, using "," to seperate them.

```go
func Router(){
  ...
  // @ApiRouter(path="/article/add", desc="add article", method="POST", group="article", handlers="add article, mid")
  article.Get("/add", articleAct.AddArticle)
  
  // @ApiRouter(path="/m/article/edit",desc="edit article",method="POST",group="article", handlers="edit article, mid")
  article.Get("/edit", articleAct.EditArticle)
  
  // @ApiRouter(path="/m/article/list",desc="list article",method="GET",group="article", handlers="list article, mid")  
  article.Get("/list", articleAct.ListArticle)
  ...
}
```

#### 控制器 handler

处理器名称必须和路由(handler)对应，否则会被抛弃.

@ApiHandler "name" has to be contained by @ApiRouter "handlers", otherwise it will be abandoned.

ApiBody 中的类型(type)被默认定义为 text, 你可以定义成 file 类型.

@ApiBody type="text" (default - text), if you don't define it, OR you could define it as 'file'.

```go
// @ApiBody(key="title", value="today news", desc="article title")

// @ApiBody(key="pic",desc="article picture",type="file",src="/eee.png")
```

```go
// @ApiHandler(name="mid")
func MiddwareHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiHeader(key="Content-Type",value="application/json",desc="header description")
  ...
}

// @ApiHandler(name="list article")
func ListHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiQuery(key="classID", value="123", desc= "article class id")
  cID := req.URL.Query().Get("classID")
  ...
}

// @ApiHandler(name="edit article")
func EditHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiQuery(key="id", value="123", desc="article id")
  ArticleID := req.URL.Query().Get("id")
  
  // @ApiBody(key="title",desc="article title")
  articleTitle := req.PostFormValue("title")
  
  // @ApiBody(key="content",desc="article content")
  articleContent := req.PostFormValue("content")
  
  // @ApiBody(key="author", value="xxx", desc="author")
  articleAuthor := req.PostFormValue("author")
  
  // @ApiBody(key="pic",desc="article picture",type="file",src="/eee.png")
  ...
}

// @ApiHandler(name="add article")
func AddHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @ApiHeader(key="Content-Type",value="application/x-www-form-urlencoded",desc="header description")
  // @ApiHeader(key="Content-Type",value="application/json",desc="header description")
  
  // @ApiQuery(key="classID", value="123", desc="article class id")
  
  // @ApiBody(key="title",desc="article title")
  // @ApiBody(key="content",desc="article content")
  // @ApiBody(key="author",desc="author")
  ...
}

```
-----

### 命令行工具 command tool

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
  -s string
      指定文件后缀名，只分析该后缀的文件，eg: ".go"
```

