package template

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"text/template/parse"
)

type treeReceiver func(tree *parse.Tree, kind Kind) error

/**
共享模板的名字, 共享模板永远不被执行.
用于优化 Clone(), 共享 common, FuncMap.
*/
const shareName = "57b11edbe4825b57ab27b7beab9848ed"

/**
从文件进行读取, 解析, 命名, hack
filename 是经过 clear 后的绝对路径.
*/
func parseFiles(t *Template, filename ...string) error {

	var (
		kind Kind
	)

	// 再次检查使用的模板名
	names := map[string]bool{}

	for _, name := range filename {

		switch path.Ext(name) {
		case "":
			return fmt.Errorf("template: invalid file name: %q", name)
		case ".html":
			kind = HTML
		default:
			kind = TEXT
		}

		// 懒方法, 直接跳过
		if t.base.Lookup(name).IsValid() {
			continue
		}

		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}

		err = parseText(t, names, kind, name, string(b))
		if err != nil {
			return err
		}

	}

	filename = []string{}

	for k, toload := range names {
		if toload {
			filename = append(filename, k)
		}
	}

	// 递归载入文件
	if len(filename) != 0 {
		return parseFiles(t, filename...)
	}

	return nil
}

/**
parseText 需要知道第一个 tree 的 kind. 以便添加到 t.
可能会载入新的文件, 产生模板名称对应的模板尚未载入.
直到全部解析完才会完整载入.
参数 ns 是给 tree 的 uri 命名.
*/
func parseText(t *Template, names map[string]bool,
	kind Kind, ns, text string) error {

	var name string

	trees, err := parse.Parse(ns, text,
		t.leftDelim, t.rightDelim, placeholderFuncs, t.base.funcs)

	if err != nil {
		return err
	}

	rootdir := t.RootDir()
	dir := absPath(ns)

	for from, tree := range trees {

		name = from
		// define 内嵌模板不能有扩展名
		if name != ns {

			if path.Ext(name) != "" {
				return fmt.Errorf(
					"template: extension are not supported on define %q", from)
			}
			name = relToURI(rootdir, dir, name)
		}

		if name == "" {
			return fmt.Errorf("template: is invalid on define %q", from)
		}

		// 需要再次检查模板是否被载入
		if t.base.Lookup(name).IsValid() {
			return fmt.Errorf(
				"template: redefinition of %q", name)
		}

		tree.Name = name
		tree.ParseName = ns

		err = hackNode(t, names, name, tree.Root, -1, tree.Root)

		if err == nil {
			_, err = t.AddParseTree(kind, tree)
		}

		if err != nil {
			return err
		}
	}
	names[ns] = false
	return nil
}

/**
hackNode 对模板中 template/import 的目标模板重新命名.
规则:
	没有扩展名当作内嵌模板,	反之当作文件模板.
	目标名称变更为绝对路径名.
	所有 template 用 import 替换.	格式为:
	import "from" "target" args...

参数:
	t     模板.
	names 所有使用的模板需要检查是否已经载入.
	from  来源模板名. 绝对路径.
	node  待 hack 的原始 parse.Node.

返回:
	是否有错误发生的错误.
*/
func hackNode(t *Template, names map[string]bool, from string,
	list *parse.ListNode, i int, node parse.Node) error {

	var (
		pipe   *parse.PipeNode
		args   []parse.Node
		target *parse.StringNode
	)
	rootdir := t.RootDir()

	switch n := node.(type) {
	default:
		return nil
	case *parse.ListNode:
		for i, node := range n.Nodes {
			err := hackNode(t, names, from, n, i, node)
			if err != nil {
				return err
			}
		}
		return nil
	case *parse.TemplateNode:

		args = make([]parse.Node, 3)

		args[0] = parse.NewIdentifier("import").SetPos(n.Pos)

		// from, 保存调用者
		args[1] = &parse.StringNode{
			NodeType: parse.NodeString,
			Pos:      n.Position(), // 伪造
			Quoted:   strconv.Quote(from),
			Text:     from,
		}

		// target, 重建目标
		args[2] = &parse.StringNode{
			NodeType: parse.NodeString,
			Pos:      n.Position(), // 伪造
			Quoted:   strconv.Quote(n.Name),
			Text:     n.Name,
		}

		// 复制其它参数
		pipe = n.Pipe
		if pipe != nil &&
			len(pipe.Cmds) != 0 &&
			pipe.Cmds[0].NodeType == parse.NodeCommand {

			for _, arg := range pipe.Cmds[0].Args {
				args = append(args, arg)
			}
		} else {
			if pipe == nil {
				pipe = &parse.PipeNode{
					NodeType: parse.NodePipe,
					Pos:      n.Position(), // 伪造
					Line:     n.Line,
					Cmds: []*parse.CommandNode{
						&parse.CommandNode{
							NodeType: parse.NodeCommand,
							Pos:      n.Position(),
						},
					},
				}
			}
		}

		pipe.Cmds[0].Args = args

		// 改成 ActionNode
		list.Nodes[i] = &parse.ActionNode{
			NodeType: parse.NodeAction,
			Pos:      n.Pos,
			Line:     n.Line,
			Pipe:     pipe,
		}

	case *parse.ActionNode:

		pipe = n.Pipe

		if pipe == nil ||
			len(pipe.Decl) != 0 ||
			len(pipe.Cmds) == 0 ||
			pipe.Cmds[0].NodeType != parse.NodeCommand ||
			len(pipe.Cmds[0].Args) == 0 ||
			pipe.Cmds[0].Args[0].Type() != parse.NodeIdentifier ||
			pipe.Cmds[0].Args[0].String() != "import" {

			return nil
		}

		args = make([]parse.Node, len(pipe.Cmds[0].Args)+1)
		args[0] = pipe.Cmds[0].Args[0]

		// from, 增加调用者来源
		args[1] = &parse.StringNode{
			NodeType: parse.NodeString,
			Pos:      args[0].Position(), // 伪造
			Quoted:   strconv.Quote(from),
			Text:     from,
		}
		// 复制其它参数
		for i, arg := range pipe.Cmds[0].Args {
			if i != 0 {
				args[i+1] = arg
			}
		}
		pipe.Cmds[0].Args = args
	}

	// 处理目标模板 args[2], 有可能是变量.
	target, _ = args[2].(*parse.StringNode)
	if target == nil {
		return nil
	}

	// 计算目标路径
	name := relToURI(rootdir, absPath(from), target.Text)
	if name == "" {
		return fmt.Errorf(
			"template: is invalid on define %q", target.Text)
	}

	// 判断文件模板是否载入
	if path.Ext(name) != "" &&
		!t.base.Lookup(name).IsValid() {

		names[name] = true
	}

	target.Text = name
	target.Quoted = strconv.Quote(name)

	return nil
}

func absPath(name string) string {

	if path.Ext(name) != "" {
		return path.Dir(name)
	}
	return path.Dir(path.Dir(name))
}

/**
absToURI 根据绝对 uri 路径 rootdir, 计算绝对路径 name 的 uri 路径.
如果返回值为 "" 表示 name 非法.
*/
func absToURI(rootdir, name string) string {

	name = cleanURI(name)
	if name == "" ||
		len(name) <= len(rootdir) ||
		name[len(rootdir)] != '/' ||
		name[:len(rootdir)] != rootdir {

		return ""
	}
	return name
}

/**
clearSlash 清理 s 中的反斜杠, 连续反斜杠, 连续斜杠.
*/
func clearSlash(s string) string {
	var to []byte

	if len(s) == 0 {
		return s
	}

	b := -1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' || c == '/' {
			if b == -1 {
				b = i
				continue
			}
			c = '/'
		}

		if b != -1 {
			if to == nil {
				// 连续斜杠
				if c == '/' && (s[i-1] == '\\' || s[i-1] == '/') {
					continue
				}

				to = make([]byte, 0, len(s))
				to = append(to, s[:b]...)
			}
			b = -1
			to = append(to, '/')
		}

		if to != nil {
			to = append(to, c)
		}
	}

	if len(to) == 0 {
		return s
	}
	return string(to)
}

/**
cleanURI 返回 uri 的最短路径写法.
*/
func cleanURI(uri string) string {

	uri = clearSlash(uri)

	if uri == "" {
		return ""
	}

	if strings.Index(uri, "/.") != -1 {
		uri = path.Clean(uri)
	}

	if uri[0] == '.' {
		return ""
	}

	return uri
}

/**
relToURI 根据绝对路径 rootdir 和 dir, 计算 name 的绝对 uri 路径.
如果返回值为 "" 表示 name 非法.
参数:
	rootdir 根目录绝对路径
	dir     当前目录绝对路径, 如果 dir 为空, 以 rootdir 替代.
	name    已经 clean 后的相对路径.
	以 "/" 开头以 rootdir 计算, 否则以 dir 计算.
*/
func relToURI(rootdir, dir, name string) string {

	if name == "" {
		return ""
	}

	if name[0] == '/' {
		return absToURI(rootdir, rootdir+name)
	}

	if dir == "" {
		return absToURI(rootdir, rootdir+"/"+name)
	}

	return absToURI(rootdir, dir+"/"+name)
}
