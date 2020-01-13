package chca

import "time"

// Article 定义了一篇文章所需要的要素
type Article struct {
	// 文章的唯一标识
	Id int `json:"id"`
	// 文章标题
	Title string `json:"title"`
	// 文章简介
	Description string `json:"description"`

	Summary  string   `json:"summary"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Category []string `json:"cate"`

	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Url       string `json:"url"`
}

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	return a[i].CreatedAt > a[j].CreatedAt
}

type Tag struct {
	Count int        `json:"count"`
	Name  string     `json:"name"`
	Posts []*Article `json:"posts"`
	Url   string     `json:"url"`
}

type Category struct {
	Count int        `json:"count"`
	Name  string     `json:"name"`
	Posts []*Article `json:"posts"`
	Url   string     `json:"url"`
}

type CollatedYear struct {
	Year   string                    `json:"year"`
	Months []*CollatedMonth          `json:"months"`
	months map[string]*CollatedMonth `json:"-"`
}

type CollatedMonth struct {
	Month string     `json:"month"`
	Posts []*Article `json:"posts"`
	month time.Month `json:"-"`
}

type Website struct {
	Title       string `yaml:"title"`
	SubTitle    string `yaml:"subtitle"`
	Description string `yaml:"description"`
	Keywords    string `yaml:"keywords"`
	SummaryLine int    `yaml:"summary_line"`
	HomeArtNum  int    `yaml:"home_art_num"`

	Theme    string `yaml:"theme"`
	Markdown string `yaml:"markdown"`
	Html     string `yaml:"html"`
	Storage  string `yaml:"storage"`

	Author string `yaml:"name"`
	Avatar string `yaml:"avatar"`
	Github string `yaml:"github"`
	Weibo  string `yaml:"weibo"`
	Mail   string `yaml:"mail"`
	Zhihu  string `yaml:"zhihu"`

	Paths []string `yaml:"paths"`
	Exts  []string `yaml:"exts"`

	UploadTheme string `yaml:"upload_theme,omitempty"`

	HomeTitle    string `yaml:"home_title,omitempty"`
	ArchiveTitle string `yaml:"archive_title,omitempty"`
	TagTitle     string `yaml:"tag_title,omitempty"`
	CateTitle    string `yaml:"cate_title,omitempty"`
	AboutTitle   string `yaml:"about_title,omitempty"`
	ArticleTitle string `yaml:"article_title,omitempty"`
}

type CollatedYears []*CollatedYear

func (c CollatedYears) Len() int {
	return len(c)
}

func (c CollatedYears) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CollatedYears) Less(i, j int) bool {
	return c[i].Year > c[j].Year
}

type CollatedMonths []*CollatedMonth

func (c CollatedMonths) Len() int {
	return len(c)
}

func (c CollatedMonths) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CollatedMonths) Less(i, j int) bool {
	return c[i].month > c[j].month
}
