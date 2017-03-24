package template

import (
	"html/template"
	"io"
	"text/template/parse"
)

type htmlTemp struct {
	*template.Template
}

/**
newHtml 新建基于 "html/template" 的共享模板.
*/
func newHtml() executor {
	return htmlTemp{
		Template: template.New(shareName),
	}
}

func (t htmlTemp) AddParseTree(tree *parse.Tree) (executor, error) {

	nt, err := t.Template.New(tree.Name).AddParseTree(tree.Name, tree)
	if err != nil {
		return nil, err
	}

	if t.Template.Tree == nil {
		t.Template = nt
		return t, nil
	}

	return htmlTemp{
		Template: nt,
	}, nil
}

func (t htmlTemp) Execute(
	p *Template, wr io.Writer, data interface{}) error {
	return t.Template.Execute(wr, data)
}

func (t htmlTemp) Funcs(funcMap FuncMap) {
	t.Template.Funcs(template.FuncMap(funcMap))
}

func (t htmlTemp) Lookup(name string) executor {

	nt := t.Template.Lookup(name)
	if nt == nil {
		return nil
	}

	return htmlTemp{
		Template: nt,
	}
}

func (t htmlTemp) Kind() Kind {
	return HTML
}

func (t htmlTemp) Tree() *parse.Tree {
	return t.Template.Tree
}
