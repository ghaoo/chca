package main

import (
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/num5/chca/conf"
	"github.com/num5/chca/template"
	"github.com/num5/chca/utils"
)

var data = map[string]interface{}{
	"sitetitle":       conf.SiteTitle(),
	"subtitle":    conf.SiteSubTitle(),
	"description": conf.SiteDescription(),
	"keywords":    conf.SiteKeywords(),
	"author":      conf.Author(),
	"avatar":      conf.Avatar(),
	"github":      conf.Github(),
	"weibo":       conf.Weibo(),
}

func Compile() {
	checkFile()
	copy()

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
}

// 编译主页
func CompileHome() {

	data["title"] = "主页"

	data["artlist"] = GetAllArt()
	data["cate"] = GetCate()
	data["tpl"] = conf.DirTheme() + "/layout/index.html"

	err := utils.MkDir(conf.DirHtml())

	if err != nil {
		panic(err)
	}

	homepath := path.Join(conf.DirHtml(), "index.html")

	htmlfile, err := os.Create(homepath)
	if err != nil {
		panic(err)
	}

	t, _ := template.New(conf.DirTheme() + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
	t.Walk(conf.DirTheme()+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译文章页
func CompileArticle() {
	artlist := GetAllArt()

	for _, art := range artlist {
		data["tpl"] = conf.DirTheme() + "/layout/post.html"

		data["title"] = art.Title
		data["description"] = art.Summary
		data["keywords"] = strings.Join(art.Tags, ",")

		data["article"] = art
		data["cate"] = GetCate()

		url := CreatePostLink(art)
		filepath := path.Join(conf.DirHtml(), url)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(conf.DirTheme() + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
		t.Walk(conf.DirTheme()+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}
}

// 编译about页
func CompileAbout() {
	about, err := GetAbout()
	if err != nil {
		panic(err)
	}

	data["title"] = "我的简历"

	data["tpl"] = conf.DirTheme() + "/layout/post.html"

	data["article"] = about
	data["cate"] = GetCate()

	filepath := path.Join(conf.DirHtml(), "about.html")

	htmlfile, err := os.Create(filepath)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(conf.DirTheme() + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
	t.Walk(conf.DirTheme()+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译归档页
func CompileArchive() {

	data["title"] = "文章归档"
	data["archive"] = GetArchive()
	data["cate"] = GetCate()
	data["tpl"] = conf.DirTheme() + "/layout/archive.html"

	filepath := path.Join(conf.DirHtml(), "archive")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(conf.DirTheme() + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
	t.Walk(conf.DirTheme()+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译cate导航页
func CompileCatePage() {

	data["title"] = "文章分类"
	data["cate"] = GetCate()
	data["tpl"] = conf.DirTheme() + "/layout/category.html"

	filepath := path.Join(conf.DirHtml(), "category")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(conf.DirTheme() + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
	t.Walk(conf.DirTheme()+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译category页面
func CompileCategory() {

	cates := GetCate()
	data["cate"] = cates

	for _, cate := range cates {

		data["title"] = "分类-" + cate.Name
		data["ptitle"] = cate.Name
		data["content"] = cate.Posts
		data["count"] = cate.Count
		data["tpl"] = conf.DirTheme() + "/layout/page.html"

		filepath := path.Join(conf.DirHtml(), "category", cate.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(conf.DirTheme() + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
		t.Walk(conf.DirTheme()+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}

}

// 编译tag导航页
func CompileTagPage() {

	data["title"] = "文章标签"
	data["cate"] = GetCate()
	data["tags"] = GetTag()
	data["tpl"] = conf.DirTheme() + "/layout/tag.html"

	filepath := path.Join(conf.DirHtml(), "tag")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, _ := template.New(conf.DirTheme() + "/layout/main.html")
	t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
	t.Walk(conf.DirTheme()+`/layout`, ".html")
	t.Execute(htmlfile, data)
}

// 编译tag页面
func CompileTag() {

	tags := GetTag()
	data["cate"] = GetCate()

	for _, tag := range tags {

		data["title"] = "标签-" + tag.Name
		data["ptitle"] = tag.Name
		data["content"] = tag.Posts
		data["count"] = tag.Count
		data["tpl"] = conf.DirTheme() + "/layout/page.html"

		filepath := path.Join(conf.DirHtml(), "tag", tag.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, _ := template.New(conf.DirTheme() + "/layout/main.html")
		t = t.Funcs(template.FuncMap{"unescaped": utils.Unescaped, "cmonth": utils.CMonth, "format": utils.Format, "count": utils.Count, "lt": utils.Lt, "gt": utils.Gt, "eq": utils.Eq})
		t.Walk(conf.DirTheme()+`/layout`, ".html")
		t.Execute(htmlfile, data)
	}

}

func CrearteMark(filename string) string {
	file := path.Join(conf.DirMark(), filename+".md")

	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		log.Println("已存在文件")
		os.Exit(1)
	}

	src, err := utils.CreateFile(conf.DirMark(), filename+".md")
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

func copy() {

	// copy 配置文件
	/*_, err := utils.CopyFile("conf.ini", path.Join(conf.DirHtml(), "conf.ini"))
	  if err != nil {
	      panic(err)
	  }*/

	err := utils.CopyDir(path.Join(conf.DirTheme(), "assets"), path.Join(conf.DirHtml(), "assets"))
	if err != nil {
		panic(err)
	}

}

func checkFile() {
	if _, err := os.Stat(conf.DirTheme()); os.IsNotExist(err) {
		log.Println("需要先初始化并添加模板文件")
		os.Exit(2)
	}

	/*if _, err := os.Stat(conf.DirStor()); os.IsNotExist(err) {
		panic("需要先初始化")
	}*/
}
