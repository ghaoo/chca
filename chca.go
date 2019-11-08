package chca

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{
		`MOD`: `chca`,
	})

	confile = "config.yml"
)

func Initialize() {

	createConf()
	createDir()
	log.Debug("初始化成功！")
}

func ListenHttpServer(port int) {

	log.Info("打开内置web服务器...")

	p := strconv.Itoa(port)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(Config().Html+"/assets/"))))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(Config().Html))))

	f := newFileHandler(Config().UploadTheme, Config().Markdown)
	f.Http()

	log.Debugf("内置web服务器开启成功，监听端口 :%d...", port)

	err := http.ListenAndServe(":"+p, nil)

	if err != nil {
		log.Errorf("ListenAndServe: %s", err)
	}
}

func init() {
	log.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	log.Logger.SetLevel(logrus.DebugLevel)
}
