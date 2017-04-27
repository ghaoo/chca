
# chca 模板变量

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

### 关于 `文章`、`分类`、`标签`、`归档` 字段说明

1. **文章**

博客中所有的文章都是从markdown解析过来的，字段包含：
- 文章标题 **Title**
- 文章描述 **Description**
- 文章概要 **Summary**
- 文章内容 **Content**
- `(字典)` 文章标签字典 **Tags**
- `(字典)` 文章分类字典 **Category**
- 文章创建时间 **CreatedAt**
- 文章链接 **Url**

2. **分类** 和 **标签**

分类和标签字段相同，包含：
- 分类或标签下文章数量 **Count**
- 分类或标签名称 **Name**
- `(字典)` 分类或标签下文章字典 **Posts**
- 分类或标签链接 **Url**

3. **归档**

文章归档按`年`、`月`做了整理：
- 发表文章的年份 `Year`
- `(字典)`该年份下每个月文章的集合 `Months`
- - 发表文章的月份 `Month`
- - `(字典)`该月份下面的所有文章的集合 `Posts`

> 字段中所有标明是`字典`的字段都需要循环出来才能使用，所有字段区分大小写

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

####  main.html 布局模板

> title 页面title

```html
<title>{{ .title }}</title>  // 页面标题，一般只用在布局模板`main.html` title里
```

> tpl 引用的其他模板页面，也就是需要展示的`主页面`

```html
<div class="content">
    {{import .tpl .}}
</div>

```

#### index.html 首页模板

> artlist 主页文章字典
 
`artlist` 里存放了显示在主页的所有文章，按时间顺序排列，文章数量配置项：`home_art_num`

artlist是存放文章的字典，需要遍历显示：
```html
{{ range .artlist }} <!--循环获取分类信息-->

    {{.Title}} <!--获取文章标题-->
    
    ... <!--其他字段-->
    
{{end}} <!--结束循环-->
```

#### post.html 文章模板

> article 页面文章

`article` 包含了文章所有字段，字段参考[`文章`](#关于-文章分类标签-字段说明)

```html
{{ .article.Title }} <!--文章标题-->
{{ .article.Content | unescaped }} <!--文章内容，unescaped函数把文章转换为 HTML 实体，见 辅助函数-->
... <!--其他你所需要的字段-->
```

#### archive.html 归档导航页

> archive 文章归档字典

#### category.html 分类导航页

> cate 文章分类字典，参考公共变量 cate

#### tag.html 标签导航页

> tags 文章标签字典，参考公共变量 tags

#### page.html 标签和分类列表页面

> ptitle 标签或分类名称
> count 标签或分类下的文章数量
> content 标签或分类下的文章字典

### 通用变量
`通用函数一般都是配置内容，可直把内容写在模板里，不用刻意使用`

----

##### avatar 头像链接地址

```html
{{ .avatar }}  // 输出 /assets/avatar.jpg，此链接为config.yml的 avatar 配置

```

##### sitetitle 网站标题

```html
{{ .sitetitle }}  // config.yml里 title 配置

```

##### subtitle 网站副标题

```html
{{ .subtitle }}  // config.yml里 subtitle 配置

```

##### description 网站mate说明

```html
{{ .description }}  // config.yml里 description 配置

```

##### keywords 网站mate关键字

```html
{{ .keywords }}  // config.yml里 keywords 配置

```

##### author 网站作者

```html
{{ .author }}  // config.yml里 author 配置

```

##### github 作者github主页地址

```html
{{ .github }}  // config.yml里 github 配置

```

##### weibo 作者微博主页地址

```html
{{ .weibo }}  // config.yml里 weibo 配置

```

##### zhihu 作者知乎主页地址

```html
{{ .zhihu }}  // config.yml里 zhihu 配置

```



### 辅助函数

> unescaped html转义

为了安全起见，在模板内容处理时会自动做escape处理，使用此函数会把内容还原为html代码
```html
    {{ .article.Content | unescaped }} 
    <!--使用 | 分隔开，竖线前面为变量值，竖线后面为辅助函数，意思是将竖线 | 左边的变量作为右边的函数的参数传递过去并输出返回值-->
```

> format

将时间戳转化为年月日，格式为 "2017-01-02"

> cmonth

将时间戳转化为月份，格式为 "01-02"

> count

计算一个字典的长度，不统计空值

> md5 

将字符串md5加密


## 参照模板

[theme](theme/blog)
