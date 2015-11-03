package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"errors"

	"github.com/guhao022/chca/conf"

)

const (
	HELP = `

Usage:

    chca command [args...]

初始化博客文件夹

    chca init

新建 markdown 文件

    chca new filename

同步 markdown 文件

    chca sync

编译博客

    chca compile

打开文件服务器

    chca http "port"

	`
)

func PrintUsage() {
	fmt.Println(HELP)
}

var (
	args      []string
)

func main() {
	flag.Parse()
	args = flag.Args()
	if len(args) == 0 || len(args) > 3 {
		PrintUsage()
		os.Exit(1)
	}

	switch args[0] {

	default:
		PrintUsage()
		os.Exit(1)
	case "init":
		new()
    case "new":
        if len(args) == 2 {
            name := args[1]

            CrearteMark(name)
        } else {
            log.Println(errors.New("缺少文件名"))
        }
    case "compile":
		Compile()
	case "http":
		var port int = 9900
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}

			port = p
		}

		_http(port)
	}
}

func new() {

	conf.InitConf()
	conf.InitDir()
}

func _http(port int) {

	p := strconv.Itoa(port)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(conf.DirHtml()+"/assets/"))))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(conf.DirHtml()))))

	log.Println("监听端口 :" + p + "...")

	err := http.ListenAndServe(":"+p, nil)

	if err != nil {
		log.Printf("ListenAndServe: %s\n", err)
	}
}