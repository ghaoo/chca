package utils

import (
	md "github.com/russross/blackfriday"
	"log"
	"regexp"
	"strings"
)

// 封装Markdown转换为Html的逻辑

var (
	TOC_TITLE = "<h4>文章导航:</h4>"
)

var navRegex = regexp.MustCompile(`(?ismU)<nav>(.*)</nav>`)

func MarkdownToHtml(content string) (str string) {
	defer func() {
		e := recover()
		if e != nil {
			str = content
			log.Println("Render Markdown ERR:", e)
		}
	}()

	htmlFlags := 0

	if strings.Contains(strings.ToLower(content), "[toc]") {

		htmlFlags |= md.HTML_TOC
	}

	htmlFlags |= md.HTML_USE_XHTML
	htmlFlags |= md.HTML_USE_SMARTYPANTS
	htmlFlags |= md.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= md.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := md.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= md.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= md.EXTENSION_TABLES
	extensions |= md.EXTENSION_FENCED_CODE
	extensions |= md.EXTENSION_AUTOLINK
	extensions |= md.EXTENSION_STRIKETHROUGH
	extensions |= md.EXTENSION_SPACE_HEADERS
	extensions |= md.EXTENSION_HARD_LINE_BREAK
	extensions |= md.EXTENSION_FOOTNOTES

	str = string(md.Markdown([]byte(content), renderer, extensions))

	if htmlFlags&md.HTML_TOC != 0 {
		found := navRegex.FindIndex([]byte(str))
		if len(found) > 0 {
			toc := str[found[0]:found[1]]
			toc = TOC_TITLE + toc
			str = str[found[1]:]
			reg := regexp.MustCompile(`\[toc\]|\[TOC\]`)
			str = reg.ReplaceAllString(str, toc)
		}
	}
	return str
}
