package utils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func NewStor(storpath, filename string) (*Storage, error) {

	// 检测文件夹是否存在   若不存在  创建文件夹
	if _, err := os.Stat(storpath); err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(storpath, os.ModePerm)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &Storage{storpath: storpath, name: filename}, nil
}

type Storage struct {
	storpath string
	name     string
}

// 将文件内容解析到value
func (sto *Storage) Get(value interface{}) error {
	var filepath = path.Join(sto.storpath, sto.name)
	return storread(filepath, value)
}

// 缓存文件
func (sto *Storage) Store(value interface{}) error {
	var filepath = path.Join(sto.storpath, sto.name)
	return storwrite(filepath, value)
}

// 删除文件
func (sto *Storage) Del() error {
	var filepath = path.Join(sto.storpath, sto.name)
	return os.Remove(filepath)
}

func Del(path string) error {
	return os.RemoveAll(path)
}

func getFile(storpath string) (*os.File, error) {
	f, err := os.OpenFile(storpath, os.O_RDWR, 0666)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return os.Create(storpath)
		}
		return nil, err
	}
	return f, nil
}

func storread(storpath string, value interface{}) error {
	f, err := os.OpenFile(storpath, os.O_RDWR, 0666)
	defer f.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(bufio.NewReader(f)).Decode(&value)
}

func storwrite(storpath string, value interface{}) error {
	content, err := json.Marshal(value)

	if err != nil {
		return err
	}
	return ioutil.WriteFile(storpath, content, os.ModePerm)
}
