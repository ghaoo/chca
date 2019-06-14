package chca

import (
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	eventTime    = make(map[string]int64)
	scheduleTime time.Time
)

type Watch struct {
	Paths []string
	Exts  []string
}

func NewWatch(paths []string, exts []string) *Watch {
	return &Watch{paths, exts}
}

func (w *Watch) Watcher() {
	//初始化监听器
	log.Info("初始化监听器... ... ")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("初始化监听器失败" + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				build := true
				if !w.checkIfWatchExt(event.Name) {
					continue
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Infof(" SKIP %s ", event)
					continue
				}

				mt := w.getFileModTime(event.Name)
				if t := eventTime[event.Name]; mt == t {
					log.Infof(" SKIP %s ", event.String())
					build = false
				}

				eventTime[event.Name] = mt

				if build {
					go func() {
						scheduleTime = time.Now().Add(1 * time.Second)
						for {
							time.Sleep(scheduleTime.Sub(time.Now()))
							if time.Now().After(scheduleTime) {
								break
							}
							return
						}
						log.Infof("触发编译事件: %s ", event)

						go Compile()
					}()
				}

			case err := <-watcher.Errors:
				log.Errorf("监控失败 %s ", err)
			}
		}
	}()

	for _, path := range w.Paths {
		log.Infof("监听文件夹: [%s] ", path)
		err = watcher.Add(path)
		if err != nil {
			log.Errorf("监视文件夹失败: [ %s ] ", err)
			os.Exit(2)
		}
	}
	log.Debug("初始化监控成功... ...")
}

func (w *Watch) checkIfWatchExt(name string) bool {
	for _, s := range w.Exts {
		if strings.HasSuffix(name, "."+s) {
			return true
		}
	}
	return false
}

func (w *Watch) getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {

		log.Errorf("文件打开失败 [ %s ]", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Errorf("获取不到文件信息 [ %s ]", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}