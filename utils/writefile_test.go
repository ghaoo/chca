package utils

import (
	"bufio"
	"os"
	"testing"
)


var content = []byte("test content\n")


func Benchmark_WirteAt(b *testing.B) {
	f, err := os.OpenFile("./WriteAt.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		b.Error("Open file error")
	}
	defer f.Close()
	var offset int64 = 0
	for i := 0; i < b.N; i++ {
		f.WriteAt(content, offset)
		offset += int64(len(content))
	}
}
func Benchmark_Wirte(b *testing.B) {
	f, err := os.OpenFile("./Write.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		b.Error("Open file error")
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Write(content)
	}
}
func Benchmark_WirteString(b *testing.B) {
	f, err := os.OpenFile("./WriteString.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		b.Error("Open file error")
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.WriteString(string(content))
	}
}
func Benchmark_BufioWrite(b *testing.B) {
	f, err := os.OpenFile("./WriteBufioWrite.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		b.Error("Open file error")
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i := 0; i < b.N; i++ {
		w.Write(content)
	}
	w.Flush()
}
func Benchmark_BufioWriteString(b *testing.B) {
	f, err := os.OpenFile("./WriteBufioWriteString.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		b.Error("Open file error")
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i := 0; i < b.N; i++ {
		w.WriteString(string(content))
	}
	w.Flush()
}
