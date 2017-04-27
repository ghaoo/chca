[toc]

# chca 模板函数

## 模板文件夹结构
- blog  模版文件夹
- -- assets  资源文件夹
- -- -- css
- -- -- js
- -- layout  布局文件夹
- -- -- main.html  公共布局
- -- -- index.html  主页
- -- -- post.html   文章页
- -- -- archive.html 导航导航页
- -- -- category.html 分类导航页
- -- -- tag.html 标签导航页
- -- -- page.html    标签、导航和分类列表页面

## 模板语法

模板使用的是标准的golang模板语法，语法教程参考：

### 公共变量

> cate 分类字典

cate里存放了网站所有分类信息，包括网站分类名称 `Name` 、文章数量 `Count` 、 文章字典 `Posts` 、 分类链接地址 `Url`

使用：
```html

{{ range .cate }} <!--循环获取分类信息-->
    <a href="{{ .Url }}"> <!--分类链接-->
        {{ .Name }} <!--分类名称-->
        {{ .Count }} <!--分类下文章数量-->
        {{range .Posts}} <!--循环获取分类下所有文章-->
            {{.Title}} <!--文章标题-->
        {{end}}
    </a>
{{ end }}

```

> tags 标签字典

tags里存放了网站所有标签信息。包括标签名称 `Name` 、文章数量 `Count` 、 文章字典 `Posts` 、 标签链接地址 `Url`，使用方法同上

### 模板页面独立变量

####  main.html

> title 页面title

使用：
```html
<title>{{ .title }}</title>  // 页面标题，一般只用在布局模板`main.html` title里
```

> tpl 引用的其他模板页面，也就是需要展示的`主页面`

使用：
```html
<div class="content">
    {{import .tpl .}}
</div>

```

#### index.html 

> artlist 主页文章字典
 
artlist里存放了显示在主页的所有文章，按时间顺序排列，文章数量配置项：`home_art_num`

使用：
```html
{{ range .artlist }} <!--循环获取文章-->

{{end}} <!--循环结束-->
```


### 通用变量
`通用函数一般都是配置内容，可直把内容写在模板里，不用刻意使用`

----

> avatar 头像链接地址

使用：
```html
{{ .avatar }}  // 输出 /assets/avatar.jpg，此链接为config.yml的 avatar 配置

```

> sitetitle 网站标题

使用：
```html
{{ .sitetitle }}  // config.yml里 title 配置

```

> subtitle 网站副标题

使用：
```html
{{ .subtitle }}  // config.yml里 subtitle 配置

```

> description 网站mate说明

使用：
```html
{{ .description }}  // config.yml里 description 配置

```

> keywords 网站mate关键字

使用：
```html
{{ .keywords }}  // config.yml里 keywords 配置

```

> author 网站作者

使用：
```html
{{ .author }}  // config.yml里 author 配置

```

> github 作者github主页地址

使用：
```html
{{ .github }}  // config.yml里 github 配置

```

> weibo 作者微博主页地址

使用：
```html
{{ .weibo }}  // config.yml里 weibo 配置

```

> zhihu 作者知乎主页地址

使用：
```html
{{ .zhihu }}  // config.yml里 zhihu 配置

```



### 辅助函数
> unescaped -- html转义

为了安全起见，在模板内容处理时会自动做escape处理，使用此函数会把内容还原为html代码
```html
```
