package action

// @pmServer{"path": "http://127.0.0.1:18080"}

// @pmRouter{"name": "添加文章","method": "Post","path": "/m/article/add","group": "文章"}

// @pmRouter{"name": "编辑文章","method": "Post","path": "/m/article/edit","group": "文章"}

// @pmRouter{"name": "文章列表","method": "Get","path": "/m/article/list","group": "文章"}

// @pmRouter{"name": "添加文章分类","method": "Post","path": "/m/article/class/add","group": "文章分类"}

// @pmRouter{"name": "编辑文章分类","method": "Post","path": "/m/article/class/edit","group": "文章分类"}

// @pmRouter{"name": "文章分类列表","method": "Get","path": "/m/article/class/list","group": "文章分类"}

// @pmRouter{"name": "首页","method": "Get","path": "/m/home"}

// @pmHandler{"name": "文章列表"}

// @pmHeader{"key":"header","value":"value","desc":"header描述"}

// @pmQuery{"key": "id","desc": "用户id"}

// @pmQuery{"key": "column","value": "id,name,age","desc": "需要的字段"}

// @pmHandler{"name":"编辑文章"}

// @pmBody{"key":"title","desc":"文章标题"}

// @pmBody{"key":"content","desc":"文章内容"}

// @pmBody{"key":"author","desc":"作者"}
