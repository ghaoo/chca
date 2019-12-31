package main

import (
	"github.com/ghaoo/chca"
	"github.com/ghaoo/chca/utils"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

func RandStr(n int) string {
	rand.Seed(time.Now().Unix())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {

	filename := "11111111"

	file := path.Join(chca.Config().Markdown, filename+".md")

	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		chca.CrearteMark(file)
	}

	for i := 1; i <= 1000; i++ {
		utils.WriteFile(file, "\n"+RandStr(RandInt(1, 30)))

		t := RandInt(10, 500)

		log.Printf("等待 %d 毫秒, 第 %d 次写入 ... ...", t, i)

		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}
