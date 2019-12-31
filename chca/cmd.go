package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/ghaoo/chca"
)

var chcaStr = `
============================================
*   _______________  __________________	   *
*   __  ____/___  / / /__  ____/___    |   *
*   _  /     __  /_/ / _  /     __  /| |   *
*   / /___   _  __  /  / /___   _  ___ |   *
*   \____/   /_/ /_/   \____/   /_/  |_|   *
*                                          *
*             Simple and fast              *
============================================


`

const (
	HELP = `

Usage:

chca command [args...]

	初始化博客文件夹
    	chca init

	新建 markdown 文件
    	chca new filename

	编译博客
    	chca compile/c

    打开文件监听器
    	chca watch/w

	打开文件服务器
    	chca http/web [port]

    运行chca所有服务，包括内置服务器、监听器
    	chca run [port]

	`
)

func PrintUsage() {
	fmt.Println(HELP)
}

var (
	args []string
)

func main() {

	color.Green(chcaStr)

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
		chca.Initialize()
	case "new":
		if len(args) == 2 {
			name := args[1]

			chca.CrearteMark(name)
		} else {
			panic("缺少文件名")
		}
	case "compile", "c":
		chca.Compile()
	case "watch", "w":
		chca.NewWatch(chca.Config().Paths, chca.Config().Exts).Watcher()
		done := make(chan bool)
		<-done
	case "run":
		chca.Compile()
		chca.NewWatch(chca.Config().Paths, chca.Config().Exts).Watcher()
		var port = 9900
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}

			port = p
		}
		chca.ListenHttpServer(port)

	case "http", "web":
		var port = 9900
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}

			port = p
		}

		chca.ListenHttpServer(port)
	}
}
