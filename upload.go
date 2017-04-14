package main

import (
	"encoding/json"
	"net/http"
	"os"
	"io"
	"fmt"
	"time"
	"path/filepath"
	"strings"
	"html/template"
)

type FileHandler struct {
	tplPath  string
	savePath string
}

func newFileHandler(tpl, save string) *FileHandler {
	return &FileHandler{
		tplPath:  tpl,
		savePath: save,
	}
}

func (fh *FileHandler) Http() {
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir(fh.tplPath+"/asset/"))))
	http.Handle("/file", http.StripPrefix("/file/", http.FileServer(http.Dir(fh.savePath))))
	http.HandleFunc("/markdown", fh.index)
	http.HandleFunc("/upload", fh.upload)
	http.HandleFunc("/files", fh.filewolk)
}

func (fh *FileHandler) index(w http.ResponseWriter, r *http.Request) {

	user := r.FormValue("u")
	password := r.FormValue("p")

	if user != "guhao" || password != "gghao" {
		w.Write([]byte("未找到页面"))
		w.WriteHeader(404)
		return
	}

	t, err := template.ParseFiles(fh.tplPath + "/index.html")
	if err != nil {
		log.Errorf("解析主页模版失败：%s", err)
	}
	err = t.Execute(w, "上传文件")
	if err != nil {
		log.Errorf("解析主页模版失败：%s", err)
	}
}

// 上传文件接口
func (fh *FileHandler) upload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("文件上传异常:%s\n", err)
		}
	}()

	if "POST" == r.Method {

		r.ParseMultipartForm(32 << 20) //在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Errorf("未找到上传文件：%s", err)
			resp := map[string]interface{} {
				"code": 500,
				"error": "未找到上传文件:"+err.Error(),
			}
			out, _ := json.Marshal(resp)
			w.Write(out)
			return
		}

		filename := handler.Filename

		save := fh.savePath + "/" + filename

		//检查文件是否存在
		if !Exist(fh.savePath) {
			os.MkdirAll(fh.savePath, os.ModePerm)
		} else {
			if Exist(save) {
				log.Warnf("博客《%s》文件已经存在", filename)
				resp := map[string]interface{} {
					"code": 500,
					"error": "博客《"+filename+"》文件已经存在",
				}
				out, _ := json.Marshal(resp)
				w.Write(out)
				return
			}
		}

		//结束文件
		of, err := handler.Open()
		if err != nil {
			log.Errorf("文件处理错误： %s", err)
			resp := map[string]interface{} {
				"code": 500,
				"error": "文件处理错误:"+err.Error(),
			}
			out, _ := json.Marshal(resp)
			w.Write(out)
			return
		}
		defer file.Close()

		//保存文件
		f, err := os.Create(save)
		if err != nil {
			log.Errorf("创建文件失败： %s", err)
			resp := map[string]interface{} {
				"code": 500,
				"error": "创建文件失败:"+err.Error(),
			}
			out, _ := json.Marshal(resp)
			w.Write(out)
			return
		}
		defer f.Close()
		io.Copy(f, of)

		//获取文件状态信息
		fstat, _ := f.Stat()

		//打印接收信息
		print := fmt.Sprintf("上传时间:%s, Size: %dKB,  Name:%s\n", time.Now().Format("2006-01-02 15:04:05"), fstat.Size()/1024, filename)
		log.Infof(print)

		resp := map[string]interface{} {
			"code": 0,
			"msg": print,
		}
		out, _ := json.Marshal(resp)
		w.Write(out)

		return
	}
}

func (fh *FileHandler) filewolk(w http.ResponseWriter, r *http.Request) {
	dir := fh.savePath
	var filemaps []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {
			return err
		}
		if f.IsDir() {
			return nil
		}

		filename := strings.TrimLeft(path, fh.savePath)

		filemaps = append(filemaps, filename)

		return nil
	})
	if err != nil {
		w.Write([]byte("filepath.Walk() returned" + err.Error()))
	}

	out, err := json.Marshal(filemaps)
	w.Write(out)
}

func (fh *FileHandler) delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filename := r.FormValue("filename")

	file := fh.savePath + "/" + filename

	err := os.Remove(file)
	if err != nil {
		resp := map[string]interface{} {
			"code": 500,
			"error": "删除文件失败:"+err.Error(),
		}
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}

	resp := map[string]interface{} {
		"code": 0,
		"error": "删除《" + filename + "》文件成功",
	}
	out, _ := json.Marshal(resp)
	w.Write(out)
	return
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
