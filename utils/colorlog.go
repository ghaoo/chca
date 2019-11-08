package utils

import (
	"fmt"
)

const (
	gray = uint8(iota + 90)
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	//NRed      = uint8(31) // Normal
)

// 输出蓝色字符
func Blue(format string, a ...interface{}) {
	fmt.Printf("\033[%dm%s\033[0m", blue, fmt.Sprintf(format, a...))
}

// 输出青色字符
func Cyan(format string, a ...interface{}) {
	fmt.Printf("\033[%dm%s\033[0m", cyan, fmt.Sprintf(format, a...))
}

// 输出红色字符
func Red(format string, a ...interface{}) {
	fmt.Printf("\033[%dm%s\033[0m", red, fmt.Sprintf(format, a...))
}

// 输出洋红色字符
func Magenta(format string, a ...interface{}) {
	fmt.Printf("\033[%dm%s\033[0m", magenta, fmt.Sprintf(format, a...))
}

// 输出绿色字符
func Green(format string, a ...interface{}) {
	fmt.Printf("\033[%dm%s\033[0m", green, fmt.Sprintf(format, a...))
}

// 输出黄色字符
func Yellow(format string, a ...interface{}) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", yellow, fmt.Sprintf(format, a...))
}
