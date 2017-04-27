package main

import (
	"os"
	"path"
	"strings"
	"time"
	
	"github.com/num5/chca/template"
	"github.com/num5/chca/utils"
)

var data = map[string]interface{}{
	"sitetitle":   Config().Title,
	"subtitle":    Config().SubTitle,
	"description": Config().Description,
	"keywords":    Config().Keywords,
	"author":      Config().Author,
	"avatar":      Config().Avatar,
	"github":      Config().Github,
	"weibo":       Config().Weibo,
	"zhihu":       Config().Zhihu,
}

func Compile() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic 错误: %s\n", err)
		}
	}()

	log.Tracf("开始编译博客...")

	checkFile()
	subcopy()

	LoadArticle()
	// 创建页面
	CompileHome()
	CompileArticle()
	CompileArchive()
	CompileTagPage()
	CompileCatePage()
	CompileCategory()
	CompileTag()
	CompileAbout()
	log.Debug("编译完成...")
}

// 编译主页
func CompileHome() {

	title := Config().HomeTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "主页"
	} else {
		data["title"] = title
	}

	data["artlist"] = GetHomeArt()
	data["cate"] = GetCate()
	data["tags"] = GetTag()
	data["tpl"] = Config().Theme + "/layout/index.html"

	err := utils.MkDir(Config().Html)

	if err != nil {
		panic(err)
	}

	homepath := path.Join(Config().Html, "index.html")

	htmlfile, err := os.Create(homepath)
	if err != nil {
		panic(err)
	}

	t, _ := template.New(Config().Theme + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
	t.Walk(Config().Theme+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译文章页
func CompileArticle() {
	artlist := GetAllArt()

	title := Config().ArticleTitle

	pfix := ""

	if len(strings.TrimSpace(title)) > 0 {
		pfix = title
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()

	for _, art := range artlist {
		data["tpl"] = Config().Theme + "/layout/post.html"

		data["title"] = pfix + "-" + art.Title
		data["description"] = strings.TrimSpace(art.Summary)
		data["keywords"] = strings.Join(art.Tags, ",")

		data["article"] = art

		url := CreatePostLink(art)
		filepath := path.Join(Config().Html, url)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(Config().Theme + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
		t.Walk(Config().Theme+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}
}

// 编译about页
func CompileAbout() {
	about, err := GetAbout()
	if err != nil {
		panic(err)
	}

	title := Config().AboutTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "我的简历"
	} else {
		data["title"] = title
	}

	data["tpl"] = Config().Theme + "/layout/post.html"

	data["article"] = about
	data["cate"] = GetCate()
	data["tags"] = GetTag()

	filepath := path.Join(Config().Html, "about.html")

	htmlfile, err := os.Create(filepath)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(Config().Theme + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
	t.Walk(Config().Theme+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译归档页
func CompileArchive() {

	title := Config().ArchiveTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章归档"
	} else {
		data["title"] = title
	}

	data["archive"] = GetArchive()
	data["cate"] = GetCate()
	data["tags"] = GetTag()
	data["tpl"] = Config().Theme + "/layout/archive.html"

	filepath := path.Join(Config().Html, "archive")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(Config().Theme + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
	t.Walk(Config().Theme+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译cate导航页
func CompileCatePage() {

	title := Config().CateTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章分类"
	} else {
		data["title"] = title
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()
	data["tpl"] = Config().Theme + "/layout/category.html"

	filepath := path.Join(Config().Html, "category")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(Config().Theme + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
	t.Walk(Config().Theme+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译category页面
func CompileCategory() {

	cates := GetCate()
	data["cate"] = cates
	data["tags"] = GetTag()

	for _, cate := range cates {

		data["title"] = "分类-" + cate.Name
		data["ptitle"] = cate.Name
		data["content"] = cate.Posts
		data["count"] = cate.Count
		data["tpl"] = Config().Theme + "/layout/page.html"

		filepath := path.Join(Config().Html, "category", cate.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(Config().Theme + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
		t.Walk(Config().Theme+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}

}

// 编译tag导航页
func CompileTagPage() {

	title := Config().TagTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章标签"
	} else {
		data["title"] = title
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()
	data["tpl"] = Config().Theme + "/layout/tag.html"

	filepath := path.Join(Config().Html, "tag")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(Config().Theme + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
	t.Walk(Config().Theme+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译tag页面
func CompileTag() {

	tags := GetTag()
	data["cate"] = GetCate()
	data["tags"] = GetTag()

	for _, tag := range tags {

		data["title"] = "标签-" + tag.Name
		data["ptitle"] = tag.Name
		data["content"] = tag.Posts
		data["count"] = tag.Count
		data["tpl"] = Config().Theme + "/layout/page.html"

		filepath := path.Join(Config().Html, "tag", tag.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(Config().Theme + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq, "md5": utils.Xmd5})
		t.Walk(Config().Theme+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}

}

func CrearteMark(filename string) string {
	file := path.Join(Config().Markdown, filename+".md")

	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		log.Errorf("已存在文件")
		os.Exit(1)
	}

	src, err := utils.CreateFile(Config().Markdown, filename+".md")
	if err != nil {
		panic(err)
	}

	date := time.Now().Format("2006-01-02")
	now := time.Now().Format("15:04:05")
	masthead := `---
date: ` + date + `
time: ` + now + `
title: ` + filename + `
categories:
-
tags:
-
-
---`
	err = utils.WriteFile(src, masthead)
	if err != nil {
		panic(err)
	}

	return src
}

func subcopy() {

	// copy 配置文件
	/*_, err := utils.CopyFile("conf.ini", path.Join(Config().Html, "conf.ini"))
	  if err != nil {
	      panic(err)
	  }*/

	err := utils.CopyDir(path.Join(Config().Theme, "assets"), path.Join(Config().Html, "assets"))
	if err != nil {
		panic(err)
	}

}

func checkFile() {
	if _, err := os.Stat(Config().Theme); os.IsNotExist(err) {
		panic("需要先初始化并添加模板文件")
	}
}
