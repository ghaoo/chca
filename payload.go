package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
    "bytes"
    "path"

    "gopkg.in/yaml.v2"
	"github.com/guhao022/chca/conf"
	"github.com/guhao022/chca/utils"
)

var (
	htmlStor = conf.DirHtml() //编译后保存的文件夹

	contents []*Article
    cates map[string]*Category
    tags map[string]*Tag
)

func LoadArticle() {

    contents = make([]*Article, 0)

    cates = make(map[string]*Category)
    tags = make(map[string]*Tag)

	mdlist := Marklist()

	for _, fi := range mdlist {
		art, err := loadContent(fi)

		if err == nil {
			art.Url = CreatePostLink(art)
			contents = append(contents,art)

            for _, _cate := range art.Category {
                cate := cates[_cate]
                if cate == nil {
                    cate = &Category{0, _cate, make([]*Article, 0), "/category/" + _cate}
                    cates[_cate] = cate
                }
                cate.Count += 1
                cate.Posts = append(cate.Posts, art)
            }

            for _, _tag := range art.Tags {
                tag := tags[_tag]
                if tag == nil {
                    tag = &Tag{0, _tag, make([]*Article, 0), "/tag/" + _tag}
                    tags[_tag] = tag
                }
                tag.Count += 1
                tag.Posts = append(tag.Posts, art)
            }

		} else {
			panic(err)
		}
	}

    sort.Sort(Articles(contents))
}

// 获取归档信息
func GetArchive() []*CollatedYear {

    collated := make(CollatedYears, 0)

    _collated := make(map[string]*CollatedYear)

    for _, post := range contents {

        year := utils.Year(post.CreatedAt)
        month := utils.Month(post.CreatedAt)
        _month := time.Unix(post.CreatedAt, 0).Month()

        yearc := _collated[year]
        if yearc == nil {
            yearc = &CollatedYear{year, make([]*CollatedMonth, 0), make(map[string]*CollatedMonth)}
            _collated[year] = yearc
        }
        monthc := yearc.months[month]
        if monthc == nil {
            monthc = &CollatedMonth{month, []*Article{}, _month}
            yearc.months[month] = monthc
        }
        monthc.Posts = append(monthc.Posts, post)
    }

    for _, yearc := range _collated {
        monthArray := make(CollatedMonths, 0)
        for _, monthc := range yearc.months {
            monthArray = append(monthArray, monthc)
        }

        sort.Sort(monthArray)

        yearc.months = nil
        yearc.Months = monthArray
        collated = append(collated, yearc)
    }

    sort.Sort(collated)
    return collated
}

//获取菜单数组
func GetCate() map[string]*Category {
    return cates
}

// 获取tag
func GetTag() map[string]*Tag {
    return tags
}

func loadContent(file string) (art *Article, err error) {

    art = &Article{}

    ctx, err := ReadMuCtx(file)

    if err != nil {
        return nil, err
    }

    sumLines := conf.SiteSumLine()

    summary, err := makeSummary(ctx.Content, sumLines)

    if err != nil {
        return nil, err
    }

    art.Title = ctx.Title
    art.Description = ctx.Description
    art.Category = ctx.Categories
    art.Tags = ctx.Tags
    art.Summary = summary
    art.Content = utils.MarkdownToHtml(ctx.Content)
    art.CreatedAt = utils.Str2Unix("2006-01-02", ctx.Date)

    return art, nil
}

// 获取所有的文章
func GetAllArt() []*Article {
	return contents
}

// 获取about内容
func GetAbout() (art *Article, err error) {
    art = &Article{}
    about := path.Join(conf.DirMark(), "/about.md")

    if _, err := os.Stat(about); os.IsNotExist(err) {
        return art, nil
    }

    content, err := ioutil.ReadFile(about)

    if err != nil {
        return nil, err
    }

    art.Title = ""
    art.Content = utils.MarkdownToHtml(string(content))
    art.CreatedAt = time.Now().Unix()

    return art, nil
}

// 获取 markdown 文件夹下所有文件
func Marklist() (mdlist []string) {
	mddir := conf.DirMark()

	filepath.Walk(mddir, func(path string, f os.FileInfo, err error) error {

		if err != nil { //忽略错误
			return err
		}

		if f.IsDir() {
			return nil
		}

		if strings.ToLower(f.Name()) == "readme.md" {
			return nil
		}

		if f.Name() == "about.md" {
			return nil
		}

		if strings.HasSuffix(f.Name(), ".md") {
			mdlist = append(mdlist, path)
		}
		return nil
	})

	return mdlist
}


// 根据文件获取摘要信息
func makeSummary(content string, lines int) (string, error) {
    buff := bufio.NewReader(bytes.NewBufferString(content))
    dst := ""
    for lines > 0 {
        line, err := buff.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }

        if strings.Contains(line, "[toc]") {
            continue
        }

        if strings.Trim(line, "\r\n\t ") == "```" {
            continue
        }

        dst += line
        lines--
    }

    return utils.MarkdownToHtml(dst), nil
}

// 根据内容获取摘要信息
func summary(content string, n int) string {
	strSlice := strings.SplitN(content, "\n", -1)

	//var summary string
	var sumSlice []string

	for i, str := range strSlice {
		if strings.Contains(str, "[toc]") {
			continue
		}

		if i >= n {
			break
		}

		sumSlice = append(sumSlice, str)
	}

	summary := strings.Join(sumSlice, "\n")

	return summary
}


// 配置生产路径
func CreatePostLink(art *Article) string {
	t := time.Unix(art.CreatedAt, 0)

	year, month, day := t.Date()

	link := fmt.Sprintf("/%s/%d/%d/%d/%s/", "article", year, month, day, utils.Convert(art.Title))

	return link
}

type mustring struct {
    Title       string
    Description string
    Date        string
    Categories  []string
    Tags        []string
    Content     string
}

func ReadMuCtx(path string) (ctx *mustring, err error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	line, err := br.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(line, "---") {
		err = errors.New("Not Start with ---   : " + path)
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	for {
		line, err = br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		if strings.HasPrefix(line, "---") {
			break
		}
		buf.WriteString(line)
	}

    err = yaml.Unmarshal(buf.Bytes(), &ctx)

    content, err := ioutil.ReadAll(br)
    if err != nil {
        return nil, err
    }

    fi, _ := f.Stat()

    if ctx.Title == "" {
        ctx.Title = strings.Replace(strings.TrimRight(fi.Name(), ".md"), conf.DirMark()+"/", "", 1)
    }

    if ctx.Date == "" {
        ctx.Date = utils.Format(fi.ModTime().Unix())
    }

    ctx.Content = string(content)

	return
}
