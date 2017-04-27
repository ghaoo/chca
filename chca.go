package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/num5/logger"
)

var (
	log *logger.Log
	confile = "config.yml"
)


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
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic 错误: %s", err)
		}
	}()

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
			log.Error("缺少文件名")
		}
	case "compile", "c":
		Compile()
	case "watch", "w":
		NewWatch(Config().Paths, Config().Exts).Watcher()
		done := make(chan bool)
		<-done
	case "run":
		Compile()
		NewWatch(Config().Paths, Config().Exts).Watcher()
		var port int = 9900
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}

			port = p
		}
		_http(port)

	case "http", "web":
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
	createConf()
	createDir()
}

func _http(port int) {

	log.Trac("打开内置web服务器...")

	p := strconv.Itoa(port)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(Config().Html+"/assets/"))))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(Config().Html))))

	f := newFileHandler(Config().UploadTheme, Config().Markdown)
	f.Http()

	log.Debugf("打开内置web服务器成功，监听端口 :%d...", port)

	err := http.ListenAndServe(":"+p, nil)

	if err != nil {
		log.Errorf("ListenAndServe: %s", err)
	}
}

func init() {

	// 初始化Log
	log = logger.NewLog(1000)
	// 设置log级别
	log.SetLevel("Debug")
	// 设置输出引擎
	log.SetEngine("file", `{"level":5, "spilt":"size", "filename":".logs/chca.log", "maxsize":15}`)
	//log.DelEngine("console")

	// 设置是否输出行号
	log.SetFuncCall(true)
}
