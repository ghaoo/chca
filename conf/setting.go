package conf

import (
	"io/ioutil"
	"os"
)

var confile = ConfigFile

var conf_setting = `[site]
title = Golune
subtitle = 十年拿大锤，看什么都是钉子
description = 专注golang，php，python，c/c++等语言开发，服务器开发
keywords = golang,php,laravel,服务器
summary_line = 10

[dir]
theme = blog
markdown = markdown
html = /data/www/golune
storage = storage

[author]
name = guhao
avatar = /assets/avatar.png
github = https://github.com/guhao022
weibo = http://weibo.com/golune`

func InitConf() {

	_, err := os.Stat(confile)
	if os.IsNotExist(err) {
		_, err := os.OpenFile(confile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		var confwrite = []byte(conf_setting)
		err = ioutil.WriteFile(confile, confwrite, 0666) //写入文件(字节数组)
		if err != nil {
			panic(err)
		}
	}
}

func InitDir() {
	_, err := os.Stat(confile)
	if os.IsNotExist(err) {
		InitConf()
	}

	_, err = os.Stat(DirHtml())
	if os.IsNotExist(err) {

		if err := os.MkdirAll(DirHtml(), os.ModePerm); err != nil {
			panic(err)
		}
	}

	_, err = os.Stat(DirMark())
	if os.IsNotExist(err) {

		if err := os.MkdirAll(DirMark(), os.ModePerm); err != nil {
			panic(err)
		}
	}

	/*_, err = os.Stat(DirStor())
	if os.IsNotExist(err) {

		if err := os.MkdirAll(DirStor(), os.ModePerm); err != nil {
			panic(err)
		}
	}*/

	_, err = os.Stat(DirTheme())
	if os.IsNotExist(err) {

		if err := os.MkdirAll("theme", os.ModePerm); err != nil {
			panic(err)
		}
	}

}
