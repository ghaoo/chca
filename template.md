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

## 模板函数

### main.html 页面


### 通用函数
通用函数可用可不用，通用函数一般都是配置内容，可直把内容写在模板里，不用刻意使用函数

----

> avatar 头像链接地址

使用：
```html
{{ .avatar }}  // 输出 /assets/avatar.jpg，此链接为config.yml的 avatar 配置

```

> sitetitle 网站标题

使用：
```html
{{ .avatar }}  // 输出 /assets/avatar.jpg，此链接为config.yml的 avatar 配置

```

### 辅助函数
> unescaped -- html转义

为了安全起见，在模板内容处理时会自动做escape处理，使用此函数会把内容还原为html代码
```html


```
