package utils

import (
	. "github.com/russross/blackfriday"
	"log"
	"regexp"
	"strings"
)

// 封装Markdown转换为Html的逻辑

const (
	//http://maruku.rubyforge.org/maruku.html#toc-generation
	TOC_MARKUP = "[toc]"
)

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

	if strings.Contains(content, TOC_MARKUP) {
		htmlFlags |= HTML_TOC
	}

	htmlFlags |= HTML_USE_XHTML
	htmlFlags |= HTML_USE_SMARTYPANTS
	htmlFlags |= HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= HTML_SMARTYPANTS_LATEX_DASHES
	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= EXTENSION_NO_INTRA_EMPHASIS
	extensions |= EXTENSION_TABLES
	extensions |= EXTENSION_FENCED_CODE
	extensions |= EXTENSION_AUTOLINK
	extensions |= EXTENSION_STRIKETHROUGH
	extensions |= EXTENSION_SPACE_HEADERS

	str = string(Markdown([]byte(content), renderer, extensions))

	if htmlFlags&HTML_TOC != 0 {
		found := navRegex.FindIndex([]byte(str))
		if len(found) > 0 {
			toc := str[found[0]:found[1]]
			toc = TOC_TITLE + toc
			str = str[found[1]:]
			str = strings.Replace(str, TOC_MARKUP, toc, -1)
		}
	}
	return str
}
