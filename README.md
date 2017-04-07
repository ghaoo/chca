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
# 网站设置
[site]
# 网站标题
title = 我的网站

# 网站次标题
subtitle = 网站标题

# 主页 mate-description 的简介
description = mate-description

# 主页 mate-keywords 关键字
keywords = mate-keywords

# 文章摘要行数（行数指的是markdown文件的行数）
summary_line = 10

# 文件夹设置
[dir]

# 网站模板
theme = blog

# 博客markdown文件存放文件夹
markdown = markdown

# 博客html文件存放文件夹
html = /var/www/html

# 数据存放文件夹，暂时未用到
storage = storage

# 作者信息设置
[author]
# 作者名称
name = my name

# 头像
avatar = /assets/avatar.png

# github 地址
github = https://github.com/num5

# 微博地址
weibo = http://weibo.com/golune

# email 地址
mail = 378999587@qq.com

```