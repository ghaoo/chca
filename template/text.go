package template

import (
	"io"
	"text/template"
	"text/template/parse"
)

type FuncMap template.FuncMap

type textTemp struct {
	*template.Template
}

/**
newText 新建基于 "text/template" 的共享模板.
*/
func newText() executor {
	return textTemp{
		Template: template.New(shareName),
	}
}

func (t textTemp) AddParseTree(tree *parse.Tree) (executor, error) {

	nt, err := t.Template.New(tree.Name).AddParseTree(tree.Name, tree)
	if err != nil {
		return nil, err
	}

	if t.Template.Tree == nil {
		t.Template = nt
		return t, nil
	}

	return textTemp{
		Template: nt,
	}, nil
}

func (t textTemp) Execute(
	p *Template, wr io.Writer, data interface{}) error {
	return t.Template.Execute(wr, data)
}

func (t textTemp) Funcs(funcMap FuncMap) {
	t.Template.Funcs(template.FuncMap(funcMap))
}

func (t textTemp) Lookup(name string) executor {

	nt := t.Template.Lookup(name)
	if nt == nil {
		return nil
	}

	return textTemp{
		Template: nt,
	}
}

func (t textTemp) Kind() Kind {
	return TEXT
}

func (t textTemp) Tree() *parse.Tree {
	return t.Template.Tree
}
