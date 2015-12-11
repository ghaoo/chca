package conf

import (
	"os"
	"path"

	"github.com/guhao022/neutron/conf"
)

var ConfigFile = "conf.ini"

const (
	site_title        = "chca"
	site_subtitle     = ""
	site_description  = ""
	site_keywords  	  = ""
	site_summary_line = 10

	dir_markdown = "markdown"
	dir_theme    = "blog"
	dir_html     = "blog"
	dir_storage  = "storage"

	author_name   = "nil"
	author_avatar = "/assets/avatar.jpg"
	author_github = "https://github.com/guhao022"
	author_weibo  = "http://weibo.com/golune"
)

func dict() conf.Dict {

	_, err := os.Stat(ConfigFile)

	if os.IsNotExist(err) {
		InitConf()
	}

	dict, err := conf.Load(ConfigFile)
	if err != nil {
		panic(err)
	}
	return dict
}

func SiteTitle() string {
	dict := dict()
	title, found := dict.GetString("site", "title")
	if !found {
		title = site_title
	}
	return title
}

func SiteSubTitle() string {
	dict := dict()
	subtitle, found := dict.GetString("site", "subtitle")
	if !found {
		subtitle = site_subtitle
	}
	return subtitle
}

func SiteKeywords() string {
	dict := dict()
	description, found := dict.GetString("site", "keywords")
	if !found {
		description = site_keywords
	}
	return description
}

func SiteDescription() string {
	dict := dict()
	description, found := dict.GetString("site", "description")
	if !found {
		description = site_description
	}
	return description
}

func SiteSumLine() int {
	dict := dict()
	line, found := dict.GetInt("site", "summary_line")
	if !found {
		line = site_summary_line
	}
	return line
}

func DirMark() string {
	dict := dict()
	mark, found := dict.GetString("dir", "markdown")
	if !found {
		mark = dir_markdown
	}
	return mark
}

func DirTheme() string {
	dict := dict()
	theme, found := dict.GetString("dir", "theme")
	if !found {
		theme = dir_theme
	}

	tpath := path.Join("theme", theme)
	return tpath
}

func DirHtml() string {
	dict := dict()
	html, found := dict.GetString("dir", "html")
	if !found {
		html = dir_html
	}
	return html
}

func DirStor() string {
	dict := dict()
	stor, found := dict.GetString("dir", "storage")
	if !found {
		stor = dir_storage
	}
	return stor
}

func Author() string {
	dict := dict()
	author, found := dict.GetString("authar", "name")
	if !found {
		author = author_name
	}
	return author
}

func Avatar() string {
	dict := dict()
	avatar, found := dict.GetString("author", "avatar")
	if !found {
		avatar = author_avatar
	}
	return avatar
}

func Github() string {
	dict := dict()
	author, found := dict.GetString("author", "github")
	if !found {
		author = author_github
	}
	return author
}

func Weibo() string {
	dict := dict()
	weibo, found := dict.GetString("author", "weibo")
	if !found {
		weibo = author_weibo
	}
	return weibo
}
