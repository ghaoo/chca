# CHCA
一个使用golang开发的博客生成器

## 安装
```bash
go get -v github.com/num5/chca
```

或者下载安装

https://github.com/num5/chca/releases

## 使用

```bash
# chca command [args...]

# 初始化博客文件夹
    chca init

# 新建 markdown 文件
    chca new filename

# 编译博客
    chca compile/c
    
# 打开文件监听器
    chca watch/w

# 打开文件服务器， 默认端口9900
    chca http [port]
    
# 运行chca， 默认端口9900
    chca run [port]
```

### chca init
> chca init 用于初始化博客，会自动生成config.yml

```yml
# config.yml

# 站点信息
title: 我的网站
subtitle: 网站标题
description: mate-description
keywords: mate-keywords
summary_line: 10 // 首页文章行数
home_art_num: 30 // 首页文章数量

# 文件夹相关
theme: theme/blog //网站模板
markdown: markdown //博客markdown文件存放文件夹
html: /data/www/html //博客html文件存放文件夹
storage: storage //数据存放文件夹，暂时未用到

# 作者信息
author: your name
avatar: /assets/avatar.png  //头像连接，一般放到assets或者avatar文件夹里
github: https://github.com/num5  //github主页
weibo: http://weibo.com/golune  //微博主页
mail: 378999587@qq.com  //email 地址
zhihu: https://www.zhihu.com/people/golune  //知乎主页
 
# 监听信息
paths:  // 监听文件夹
  - markdown
exts:  // 监听后缀名
  - md

# 上传信息
upload_theme: theme/upload  // 上传模版地址

# 自定义标题，可不配置，使用chca设置好的标题
home_title: 主页标题
archive_title: 文章归档标题
tag_title: 标签导航页面标题
cate_title: 分类导航页面标题
about_title: 简历页面标题
article_title: 文章标题标头

```
初始化以后需要在config.yml文件同目录下创建theme文件夹用于存放模板文件

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


### chca new filename

> 新建markdown文件

markdown文件需要以 `---` 开头进行说明：

```bash

---
date: 2017-01-01
title: 我的博客
categories:
- 技术
tags:
- golang
---
```
建议使用chca创建markdown文件

about.md文件是存放作者简历的文件，存放在markdown文件夹
readme.md和about.md不会被文章解析器解析

### chca compile

> 生成html文件

### chca watch/w

> 开启文件监听器，监听文件夹和文件后缀名在config.yml里配置，配置示例：
  
  ```yml
  // 监听文件夹
  paths:
    - markdown
   
  // 监听后缀名
  exts:
    - md  // 监听 markdown 文件（以.md为后缀的文件）
  ```
  以上配置：监听器监听`markdown`文件夹下的以`.md`为后缀的文件，当文件夹下的`.md`文件新增或者发生改变时，chca则会自动编译博客

### chca http 8800

> 打开内部服务器，监听端口8800

### chca run 8800

> 打开内部服务器，监听端口8800，并开启文件监听器

# License

etcd is under the Apache 2.0 license. See the LICENSE file for details.


