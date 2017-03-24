package utils

import "html/template"

func Unescaped(x string) interface{} {
	return template.HTML(x)
}
