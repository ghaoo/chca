package chca

import (
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	//"github.com/ghaoo/chca/template"
	"github.com/ghaoo/chca/utils"
)

var data = map[string]interface{}{
	"title":       conf.Title,
	"subtitle":    conf.SubTitle,
	"description": conf.Description,
	"keywords":    conf.Keywords,
	"author":      conf.Author,
	"avatar":      conf.Avatar,
	"github":      conf.Github,
	"weibo":       conf.Weibo,
	"zhihu":       conf.Zhihu,
}

var funcMaps = template.FuncMap{
	"unescaped": utils.Unescaped,
	"cmonth":    utils.CMonth,
	"format":    utils.Format,
	"count":     utils.Count,
	"lt":        utils.Lt,
	"gt":        utils.Gt,
	"eq":        utils.Eq,
	"md5":       utils.Xmd5,
}

func Compile() {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic recovered from: %v", r)
		}
	}()

	log.Info("开始编译博客...")

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
	storageBlogMap()
	log.Debug("编译完成...")
}

func storageBlogMap() {
	stor, err := utils.NewStor(conf.Storage, "blogmap.json")
	if err != nil {
		panic(err)
	}

	err = stor.Store(GetAllArt())
	if err != nil {
		panic(err)
	}

}

// 编译主页
func CompileHome() {
	title := conf.HomeTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "主页"
	} else {
		data["title"] = title
	}

	data["artlist"] = GetHomeArt()
	data["cate"] = GetCate()
	data["tags"] = GetTag()

	err := utils.MkDir(conf.Html)

	if err != nil {
		panic(err)
	}

	homepath := path.Join(conf.Html, "index.html")

	htmlfile, err := os.Create(homepath)
	if err != nil {
		panic(err)
	}

	t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/main.tpl", conf.Theme+"/layout/home.tpl")
	if err != nil {
		panic(err)
	}

	err = t.Execute(htmlfile, data)
	if err != nil {
		panic(err)
	}
}

// 编译文章页
func CompileArticle() {
	artlist := GetAllArt()

	title := conf.ArticleTitle

	pfix := ""

	if len(strings.TrimSpace(title)) > 0 {
		pfix = title + "-"
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()

	for _, art := range artlist {

		data["title"] = pfix + art.Title
		data["description"] = strings.TrimSpace(art.Summary)
		data["keywords"] = strings.Join(art.Tags, ",")

		data["article"] = art

		url := CreatePostLink(art)
		filepath := path.Join(conf.Html, url)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/post.tpl", conf.Theme+"/layout/main.tpl")
		if err != nil {
			panic(err)
		}
		err = t.Execute(htmlfile, data)
		if err != nil {
			panic(err)
		}
	}
}

// 编译about页
func CompileAbout() {
	about, err := GetAbout()
	if err != nil {
		panic(err)
	}

	title := conf.AboutTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "我的简历"
	} else {
		data["title"] = title
	}

	data["article"] = about
	data["cate"] = GetCate()
	data["tags"] = GetTag()

	filepath := path.Join(conf.Html, "about.html")

	htmlfile, err := os.Create(filepath)

	if err != nil {
		panic(err)
	}

	t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/post.tpl", conf.Theme+"/layout/main.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(htmlfile, data)
	if err != nil {
		panic(err)
	}
}

// 编译归档页
func CompileArchive() {

	title := conf.ArchiveTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章归档"
	} else {
		data["title"] = title
	}

	data["archive"] = GetArchive()
	data["cate"] = GetCate()
	data["tags"] = GetTag()

	filepath := path.Join(conf.Html, "archive")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/archive.tpl", conf.Theme+"/layout/main.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(htmlfile, data)
	if err != nil {
		panic(err)
	}
}

// 编译cate导航页
func CompileCatePage() {

	title := conf.CateTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章分类"
	} else {
		data["title"] = title
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()

	filepath := path.Join(conf.Html, "category")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/category.tpl", conf.Theme+"/layout/main.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(htmlfile, data)
	if err != nil {
		panic(err)
	}
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

		filepath := path.Join(conf.Html, "category", cate.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/page.tpl", conf.Theme+"/layout/main.tpl")
		if err != nil {
			panic(err)
		}
		err = t.Execute(htmlfile, data)
		if err != nil {
			panic(err)
		}
	}

}

// 编译tag导航页
func CompileTagPage() {

	title := conf.TagTitle

	if len(strings.TrimSpace(title)) == 0 {
		data["title"] = "文章标签"
	} else {
		data["title"] = title
	}

	data["cate"] = GetCate()
	data["tags"] = GetTag()

	filepath := path.Join(conf.Html, "tag")

	err := utils.MkDir(filepath)

	if err != nil {
		panic(err)
	}

	filename := path.Join(filepath, "index.html")

	htmlfile, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/tag.tpl", conf.Theme+"/layout/main.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(htmlfile, data)
	if err != nil {
		panic(err)
	}
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

		data["tpl"] = conf.Theme + "/layout/page.html"

		filepath := path.Join(conf.Html, "tag", tag.Name)

		err := utils.MkDir(filepath)

		if err != nil {
			panic(err)
		}

		filename := path.Join(filepath, "index.html")

		htmlfile, err := os.Create(filename)

		if err != nil {
			panic(err)
		}

		t, err := template.New("main.tpl").Funcs(funcMaps).ParseFiles(conf.Theme+"/layout/page.tpl", conf.Theme+"/layout/main.tpl")
		if err != nil {
			panic(err)
		}
		err = t.Execute(htmlfile, data)
		if err != nil {
			panic(err)
		}
	}

}

func CrearteMark(filename string) string {
	file := path.Join(conf.Markdown, filename+".md")

	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		log.Errorf("已存在文件")
		os.Exit(1)
	}

	src, err := utils.CreateFile(conf.Markdown, filename+".md")
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

	err := utils.CopyDir(path.Join(conf.Theme, "assets"), path.Join(conf.Html, "assets"))
	if err != nil {
		panic(err)
	}

}

func checkFile() {
	if _, err := os.Stat(conf.Theme); os.IsNotExist(err) {
		panic("需要先初始化并添加模板文件")
	}
}
