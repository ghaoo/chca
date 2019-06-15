package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func BasePath() string {
	file, _ := exec.LookPath(os.Args[0])

	fpath, _ := filepath.Abs(file)

	basePath := path.Dir(fpath)

	return basePath
}

func CreateFile(dir string, name string) (string, error) {
	src := path.Join(dir, name)

	_, err := os.Stat(src)

	if os.IsExist(err) {
		return src, nil
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		if os.IsPermission(err) {
			panic("你不够权限创建文件")
		}
		return "", err
	}

	_, err = os.Create(src)
	if err != nil {
		return "", err
	}

	return src, nil
}

func MkDir(filepath string) error {

	if _, err := os.Stat(filepath); err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(filepath, os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("Source is not a directory")
	}

	// ensure dest dir does not already exist
	_, err = os.Open(dest)
	if os.IsExist(err) {

		err = os.RemoveAll(dest)
		if err != nil {
			return err
		}
		//return errors.New("Destination already exists")
	}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				panic(err)
			}
		} else {
			// perform copy
			_, err = CopyFile(sfp, dfp)
			if err != nil {
				panic(err)
			}
		}

	}
	return
}

func WriteFile(file string, text string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Errorf("Open file error: %s", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.Write([]byte(text))
	if err != nil {
		return err
	}
	return w.Flush()
}
