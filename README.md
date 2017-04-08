# CHCA
一个使用golang开发的博客生成器

### 安装
```bash
go get -v github.com/num5/chca
```

### 使用

```bash
# chca command [args...]

# 初始化博客文件夹
    chca init

# 新建 markdown 文件
    chca new filename

#编译博客
    chca compile

# 打开文件服务器，必要参数 port
    chca http "port"
```

#### chca init
> chca init 用于初始化博客，会自动生成conf.ini

```go
# conf.ini

[site]  //网站设置

title = 我的网站  //网站标题
subtitle = 网站标题  //网站次标题
description = mate-description  //主页 mate-description 的简介
keywords = mate-keywords  //主页 mate-keywords 关键字
summary_line = 10  //文章摘要行数（行数指的是markdown文件的行数）

[dir] //文件夹设置

theme = blog //网站模板
markdown = markdown  //博客markdown文件存放文件夹
html = /var/www/html  //博客html文件存放文件夹
storage = storage  //数据存放文件夹，暂时未用到

[author] //作者信息设置

name = my name  //作者名称
avatar = /assets/avatar.png  //头像
github = https://github.com/num5  //github 地址
weibo = http://weibo.com/golune  //微博地址
mail = 378999587@qq.com  //email 地址

```
初始化以后需要在conf.ini文件同目录下创建theme文件夹用于存放模板文件

> **关于模板文件的创建**

模板文件夹结构
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

关于 `标签页面` ，可以从关键字


