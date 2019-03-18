## go-postman

分析源代码中的注释，解析成postman的json文件，方便导入postman。
swagger风格。

### 注释格式
```go
服务路径
// @pmServer(path = "http://127.0.0.1:18080")
```
路由
```go
路由组
// @pmRouter(name="添加文章", method="Post", path="/m/article/add", group="文章")
// @pmRouter(name= "编辑文章",method= "Post",path= "/m/article/edit",group= "文章")
// @pmRouter(name= "文章列表",method= "Get",path= "/m/article/list",group= "文章")
单独路由
// @pmRouter(name= "首页",method= "Get",path= "/m/home")
```
处理器
```go
处理器，名字必须对应路由的名字，否则会被丢弃
// @pmHandler(name= "文章列表")
处理器中的参数
// @pmHeader(key="header",value="value",desc="header描述")
// @pmQuery(key= "id",desc= "用户id")
// @pmQuery(key="column",value= "id,name,age",desc= "需要的字段")
```
处理器
```go
处理器，名字必须对应路由的名字，否则会被丢弃
// @pmHandler(name="编辑文章")
处理器中的参数
// @pmBody(key="title",desc="文章标题")
// @pmBody(key="content",desc="文章内容")
// @pmBody(key="author",desc="作者")
```
处理器
```go
处理器，名字必须对应路由的名字，否则会被丢弃
// @pmHandler(name="添加文章")
处理器中的参数
// @pmBody(key="title",desc="文章标题")
// @pmBody(key="content",desc="文章内容")
// @pmBody(key="author",desc="作者")
// @pmBody(key="author",desc="图片",type="file",src="/eee.png")
```
