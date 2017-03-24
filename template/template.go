package template

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template/parse"
)

var (
	errUnavailable   = errors.New("template: unavailable action.")
	errUnimplemented = errors.New("template: unimplemented action.")
	errImport        = errors.New("template: unimplemented import.")
	errMuseBeNamed   = errors.New("template: must be named before.")
	errNeverShow     = errors.New("template: never show")
)

var (
	// 占位用
	placeholderFuncs = FuncMap{
		"import": func(from, name string, data ...interface{}) (template.HTML, error) {
			return "", nil
		},
	}
)

/**
模板风格
*/
type Kind uint8

const (
	INVALID Kind = iota // 无效模板
	TEXT                // 文本模板, 使用 "text/template" 进行处理.
	HTML                // html模板, 使用 "html/template" 进行处理.
)

func musetBeOverWriteThisMethod() {
	panic("must be overwrite this method.")
}
func unExpected() {
	panic("Unexpected")
}

type executor interface {
	AddParseTree(tree *parse.Tree) (executor, error)
	Execute(p *Template, wr io.Writer, data interface{}) error
	Funcs(funcMap FuncMap)
	Lookup(name string) executor
	Kind() Kind
	Name() string
	Tree() *parse.Tree
}

/**
kindTree 对 parse.Tree 进行包装, 包含 Kind 信息.
*/
type kindTree struct {
	*parse.Tree
	Kind
}

/**
Copy 返回一份 kindTree 的拷贝.
*/
func (t kindTree) Copy() kindTree {

	var tree *parse.Tree
	if t.Tree != nil {
		tree = t.Tree.Copy()
	}

	return kindTree{
		Tree: tree,
		Kind: t.Kind,
	}
}

/**
Dir 返回 ParseName 所属的目录.
*/
func (t kindTree) Dir() string {
	if t.Tree == nil {
		return ""
	}
	return path.Dir(t.Tree.ParseName)
}

/**
IsValid 返回 t 是否有效.
*/
func (t kindTree) IsValid() bool {
	return (t.Kind == TEXT || t.Kind == HTML) &&
		t.Tree != nil && t.Tree.Name != "" && t.Tree.ParseName != ""
}

/**
Name 返回 t.Tree.Name, 如果 t.Tree 为 nil, 返回 "".
*/
func (t kindTree) Name() string {
	if t.Tree == nil {
		return ""
	}
	return t.Tree.Name
}

/**
ParseName 返回 t.Tree.ParseName, 如果 t.Tree 为 nil, 返回 "".
*/
func (t kindTree) ParseName() string {
	if t.Tree == nil {
		return ""
	}
	return t.Tree.ParseName
}

/**
base 维护 Template 的基础数据. 包括: rootdir, Tree, FuncMap.
注意 Tree, FuncMap 都是用了 map 做容器, 并发读写是不安全的.
*/
type base struct {
	rootdir string
	trees   map[string]kindTree
	funcs   FuncMap
}

/**
newBase 新建 base 对象.
参数:
	rootdir 是已经处理好的 uri 格式.
*/
func newBase(rootdir string) *base {
	c := &base{
		rootdir: rootdir,
		trees:   make(map[string]kindTree),
		funcs:   make(FuncMap),
	}
	return c
}

/**
AddTree 增加 Tree.
Tree 是已经处理好的, name 要有 rootdir 前缀.
返回:
	失败返回错误信息, 成功返回 nil.
*/
func (b *base) AddTree(tree kindTree) error {

	if !tree.IsValid() {
		return fmt.Errorf("template: invalid Tree from base.AddTree.")
	}

	if len(b.rootdir) != 0 {
		if len(tree.Tree.Name) < len(b.rootdir) ||
			len(tree.Tree.ParseName) < len(b.rootdir) ||
			tree.Tree.Name[:len(b.rootdir)] != b.rootdir ||
			tree.Tree.ParseName[:len(b.rootdir)] != b.rootdir ||
			tree.Tree.Name[len(b.rootdir)] != '/' ||
			tree.Tree.ParseName[len(b.rootdir)] != '/' {

			return fmt.Errorf(
				"template: %q not under %q.", tree.Name, b.rootdir)
		}
	}

	if _, exist := b.trees[tree.Tree.Name]; exist {
		return fmt.Errorf(
			"template: redefinition of %q", tree.Tree.Name)
	}

	b.trees[tree.Tree.Name] = tree
	return nil
}

/**
Copy 返回一份 *base 的拷贝.
通常, 如果不需要即时使用 AddTree, Funcs 无需使用 Copy,
如果需要, 应保留一份 base 作为原本专用于 Copy.
*/
func (b *base) Copy() *base {

	trees := make(map[string]kindTree)

	for name, tree := range b.trees {
		trees[name] = tree.Copy()
	}

	funcs := make(map[string]interface{}, len(b.funcs))
	for k, a := range b.funcs {
		funcs[k] = a
	}

	return &base{
		rootdir: b.rootdir,
		trees:   trees,
		funcs:   funcs,
	}
}

/**
Funcs 设置自定义 FuncMap.
*/
func (b *base) Funcs(funcMap FuncMap) {
	for k, i := range funcMap {
		b.funcs[k] = i
	}
}

/**
Lookup 返回 name 对应的 kindTree. 需要使用者判断 kindTree 是否有效.
*/
func (b *base) Lookup(name string) kindTree {
	return b.trees[name]
}

/**
RootDir 返回 rootdir.
*/
func (b *base) RootDir() string {
	return b.rootdir
}

/**
Template.
*/
type Template struct {
	base       *base
	executor            // 当前执行器
	text       executor // text 执行器
	html       executor // html 执行器
	leftDelim  string
	rightDelim string
}

func (t *Template) wrap(exec executor) *Template {
	if exec == nil || exec.Kind() == INVALID {
		return nil
	}
	nt := &Template{
		base:       t.base,
		text:       t.text,
		html:       t.html,
		leftDelim:  t.leftDelim,
		rightDelim: t.rightDelim,
	}
	switch exec.Kind() {
	default:
		return nil
	case TEXT:
		nt.executor = exec
	case HTML:
		nt.executor = exec
	}
	return t
}

var dataNil = []interface{}{nil}

func (t *Template) initFuncs() *Template {

	rootdir := t.RootDir()

	funcs := map[string]interface{}{
		"import": func(from, name string,
			data ...interface{}) (template.HTML, error) {

			var (
				buf  bytes.Buffer
				exec executor
			)

			// 先假设为绝对路径, 找不到, 再求相对路径.
			exec = t.lookup(name)
			if exec == nil {
				exec = t.lookup(relToURI(
					rootdir, t.base.trees[from].Dir(), cleanURI(name)))
			}

			if exec == nil {
				return "", fmt.Errorf("template: %q is undefined", name)
			}

			if len(data) == 0 {
				data = dataNil
			}

			err := exec.Execute(t, &buf, data[0])
			if err != nil {
				return "", err
			}
			return template.HTML(buf.String()), nil
		},
	}

	t.text.Funcs(funcs)
	t.html.Funcs(funcs)

	return t
}

/**
New 基于资源路径 uri 新建一个 Template.
先对 uri 进行绝对路径计算, 计算出 rootdir 和是否要加载文件.

参数:
	uri 资源路径可以是目录或者文件, 无扩展名当作目录, 否则当作文件.
		如果 uri 为空, 用 os.Getwd() 获取目录.
		如果 uri 以 `./` 或 `.\` 开头自动加当前路径, 否则当作绝对路径.
		如果 uri 含扩展当作模板文件, 使用 ParseFiles 解析.
		uri 所指的目录被设置为 rootdir, 后续载入的文件被限制在此目录下.

	funcMap 可选自定义 FuncMap.
		当 uri 为文件时, funcMap 参数可保障正确解析模板中的函数.
返回:
	模板实例和发生的错误.
*/
func New(uri string, funcMap ...FuncMap) (*Template, error) {

	var err error

	if uri == "" {
		uri, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	} else if len(uri) > 1 && (uri[:2] == `./` ||
		uri[:2] == `.\`) {

		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		uri = dir + `/` + uri
	}

	rootdir := cleanURI(uri)
	if rootdir == "" {
		return nil, fmt.Errorf("template: invalid uri: %q", uri)
	}

	ext := path.Ext(rootdir)

	t := &Template{
		text: newText(),
		html: newHtml(),
	}

	if ext == "" {
		t.base = newBase(rootdir)
	} else {
		t.base = newBase(path.Dir(rootdir))
	}

	t.initFuncs() // init import

	for _, funcs := range funcMap {
		t.Funcs(funcs)
	}

	if ext != "" {
		err = t.ParseFiles(rootdir)
	}

	if err != nil {
		return nil, err
	}
	return t, nil
}

/**
AddParseTree 添加 tree.
参数:
	kind 值为TEXT 或 HTML, 指示 tree 采用何种风格执行.
	tree 是已经处理好的, 且 tree.Name 对应绝对路径的模板名称.
返回:
	如果 tree 符合要求, 返回 tree 对应的 *Template.
	否则返回 nil 和错误.

细节:
	事实上 ParseFiles, ParseGlob 都调用了 AddParseTree.
	如果 t 没有对应的执行模板, 自动绑定第一个 Tree 对应的模板.
*/
func (t *Template) AddParseTree(
	kind Kind, tree *parse.Tree) (*Template, error) {

	var e executor

	err := t.base.AddTree(kindTree{
		Tree: tree,
		Kind: kind,
	})

	if err == nil {
		if kind == TEXT {
			e, err = t.text.AddParseTree(tree)
		} else {
			e, err = t.html.AddParseTree(tree)
		}
	}

	if err != nil {
		delete(t.base.trees, tree.Name)
		return nil, err
	}

	if t.executor == nil {
		t.executor = e
		return t, nil
	}

	return t.wrap(e), nil
}

/**
Copy 返回一份 *Template 的拷贝. 这是真正的拷贝.
非并发安全,	如果需要 Copy 功能, 应保留一份母本专用于 Copy.
提示:
	Copy 会重建 FuncMap 中的 "import" 函数.
*/
func (t *Template) Copy() (*Template, error) {

	nt := &Template{
		base:       t.base.Copy(),
		text:       newText(),
		html:       newHtml(),
		leftDelim:  t.leftDelim,
		rightDelim: t.rightDelim,
	}

	// FuncMap
	nt.initFuncs()
	nt.text.Funcs(nt.base.funcs)
	nt.html.Funcs(nt.base.funcs)

	// 保持当前 executor 类型
	tree := t.base.Lookup(t.executor.Name())
	name := ""
	switch tree.Kind {
	default:
		nt.executor = nt.text

	case TEXT:
		name = tree.Tree.Name
		nt.executor = nt.text
		nt.executor.AddParseTree(tree.Tree)

	case HTML:
		name = tree.Tree.Name
		nt.executor = nt.html
		nt.executor.AddParseTree(tree.Tree)
	}

	// 重建环境
	for k, tree := range t.base.trees {
		if k == name {
			continue
		}
		switch tree.Kind {
		case TEXT:
			nt.text.AddParseTree(tree.Tree)
		case HTML:
			nt.html.AddParseTree(tree.Tree)
		}
	}

	return nt, nil
}

/**
Execute 执行模板, 把结果写入 wr.
*/
func (t *Template) Execute(
	wr io.Writer, data interface{}) error {

	return t.executor.Execute(t, wr, data)
}

/**
ExecuteTemplate 执行 name 对应的模板, 把结果写入 wr.
此方法先调用 Lookup 获取 name 对应的模板, 然后执行它.
*/
func (t *Template) ExecuteTemplate(
	wr io.Writer, name string, data interface{}) error {

	a := t.Lookup(name)
	if a == nil {
		return fmt.Errorf("template: %q is undefined", name)
	}
	return a.Execute(wr, data)
}

/**
Delims 设置模板定界符. 返回 t.
*/
func (t *Template) Delims(left, right string) *Template {
	t.leftDelim, t.rightDelim = left, right
	return t
}

/**
Dir 返回 t 所在目录绝对路径. slash 分割. 尾部没有 slash.
*/
func (t *Template) Dir() string {
	ns := t.executor.Name()
	if ns == shareName {
		return t.base.rootdir
	}

	if path.Ext(ns) != "" {
		return path.Dir(ns)
	}
	return path.Dir(path.Dir(ns))
}

/**
Funcs 给模板绑定自定义 FuncMap.
参数:
	funcMap 设定一次, 在所有相关模板中都会生效.
返回: t
*/
func (t *Template) Funcs(funcMap FuncMap) *Template {
	t.base.Funcs(funcMap)
	t.text.Funcs(funcMap)
	t.html.Funcs(funcMap)
	return t
}

/**
Lookup 取出 name 对应的 *Template.
参数:
	name 模板名, 相对路径. 如果以 "/" 开头表示从 rootdir 开始,
	否则从 t.Dir() 所在目录开始.
返回:
	返回 name 对应模板, 如果 name 为空或者未找到对应模板, 返回 nil.
*/
func (t *Template) Lookup(name string) *Template {

	// 计算绝对路径
	name = relToURI(t.base.rootdir, t.Dir(), cleanURI(name))
	if name == "" {
		return nil
	}

	return t.wrap(t.lookup(name))
}

func (t *Template) lookup(name string) executor {

	var exec, first, second executor

	ext := path.Ext(name)

	// 内嵌模板有可能没有扩展名
	if ext == ".html" || t.Kind() == HTML {
		first, second = t.html, t.text
	} else {
		first, second = t.text, t.html
	}

	exec = first.Lookup(name)
	if exec == nil {
		exec = second.Lookup(name)
	}

	return exec
}

/**
Name 返回 uri 风格的模板名, 事实是模板对应的绝对路径.
如果为空表示模板无效.
*/
func (t *Template) Name() string {

	if t.executor == nil {
		return ""
	}
	s := t.executor.Name()
	if s == shareName {
		return ""
	}
	return s
}

/**
ParseFiles 解析多个模板文件. 自动跳过重复的文件.
参数:
	filename 模板文件, 可使用相对路径或绝对路径.
返回:
	是否有错误发生.
*/
func (t *Template) ParseFiles(
	filename ...string) error {

	var name string
	rootdir := t.RootDir()

	for i := 0; i < len(filename); i++ {
		name = absToURI(rootdir, filename[i])
		if name == "" || path.Ext(name) == "" {
			return fmt.Errorf(
				"template: invalid filename: %q", filename[i])
		}

		filename[i] = name
	}

	err := parseFiles(t, filename...)

	if err != nil {
		return err
	}

	return nil
}

/**
Parse 解析模板源代码 text, 并以 name 命名解析后的模板.
参数:
	name 模板名字, 相对于 rootdir 的绝对路径名.
	text 待解析的模板源代码.
返回:
	解析后的模板和发生的错误.
*/
func (t *Template) Parse(name, text string) (*Template, error) {

	ns := relToURI(t.RootDir(), t.Dir(), cleanURI(name))

	if ns == "" {
		return nil, fmt.Errorf(
			"template: invalid name %q", name)
	}

	kind := TEXT
	if path.Ext(ns) == ".html" {
		kind = HTML
	}

	// 再次检查使用的模板名
	names := map[string]bool{}

	err := parseText(t, names, kind, ns, text)

	if err != nil {
		return nil, err
	}

	filename := []string{}

	for k, toload := range names {
		if toload {
			filename = append(filename, k)
		}
	}

	// 递归载入文件
	if len(filename) != 0 {
		err = parseFiles(t, filename...)
	}

	if err != nil {
		return nil, err
	}

	return t.Lookup(ns), nil
}

/**
ParseGlob 解析多个模板文件.
自动跳过重复的文件.
参数:
	pattern 模板文件模式匹配.
返回:
	是否有错误发生.
*/
func (t *Template) ParseGlob(pattern string) error {
	filename, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	fmt.Println(filename)
	return t.ParseFiles(filename...)
}

/**
RootDir 返回 rootdir.
*/
func (t *Template) RootDir() string {
	return t.base.RootDir()
}

/**
Walk 遍历 dir, 根据允许的扩展名加载模板文件.
要求所有加载文件必须位于 rootdir 之下. 自动跳过重复的文件.
参数:
	dir  要遍历的目录.
	exts 允许的扩展名拼接字符串, 格式实例: ".html.tmpl".
*/
func (t *Template) Walk(dir string, exts string) error {

	filename := []string{}
	filepath.Walk(dir,
		func(name string, info os.FileInfo, err error) error {

			if err != nil || info.IsDir() {
				return nil
			}

			name = clearSlash(name)
			ext := path.Ext(name)
			if ext == "" {
				return nil
			}

			pos := strings.Index(exts, ext) + len(ext)
			if pos >= len(ext) && (pos == len(exts) || exts[pos] == '.') {
				filename = append(filename, name)
			}

			return nil
		})

	return t.ParseFiles(filename...)
}
