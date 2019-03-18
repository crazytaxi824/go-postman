## go-postman

分析源代码中的注释，解析成postman的json文件，方便导入postman。
swagger风格。

### 事例

服务路径
```go
func server(){
  // @pmServer(path = "http://127.0.0.1:18080")
  srv.ListenAndServe := ...
}
```
路由
```go
func Router(){
  ...
  // @pmRouter(name="添加文章", method="Post", path="/article/add", group="文章")
  article.Get("/add", articleAct.AddArticle)
  
  // @pmRouter(name= "编辑文章",method= "Post",path= "/m/article/edit",group= "文章")
  article.Get("/edit", articleAct.EditArticle)
  
  // @pmRouter(name= "文章列表",method= "Get",path= "/m/article/list",group= "文章")  
  article.Get("/list", articleAct.ListArticle)
  ...
}
```
控制器
```go
// @pmRouter(name= "文章列表",method= "Get",path= "/m/article/list",group= "文章")  
func ListHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @pmQuery(key= "classID", value="123", desc= "文章类型id")
  cID := req.URL.Query().Get("classID")
  ...
}

// @pmRouter(name= "编辑文章",method= "Post",path= "/m/article/edit",group= "文章")
func EditHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @pmQuery(key= "id", value="123", desc= "文章id")
  ArticleID := req.URL.Query().Get("id")
  
  // @pmBody(key="title",desc="文章标题")
  articleTitle := req.PostFormValue("title")
  
  // @pmBody(key="content",desc="文章内容")
  articleContent := req.PostFormValue("content")
  
  // @pmBody(key="author", value="xxx", desc="作者")
  articleAuthor := req.PostFormValue("author")
  
  // @pmBody(key="author",desc="图片",type="file",src="/eee.png")
  ...
}

// @pmRouter(name="添加文章", method="Post", path="/article/add", group="文章")
func AddHandler(w http.ResponseWriter, req *http.Request) {
  ...
  // @pmHeader(key="Content-Type",value="application/x-www-form-urlencoded",desc="header描述")
  // @pmHeader(key="Content-Type",value="application/json",desc="header描述")
  
  // @pmQuery(key= "classID", value="123", desc= "文章类型id")
  
  // @pmBody(key="title",desc="文章标题")
  // @pmBody(key="content",desc="文章内容")
  // @pmBody(key="author",desc="作者")
  ...
}

```


