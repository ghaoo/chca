package utils
import (
    "strings"
)

func Convert(str string) string {
    str = strings.ToLower(str)
    ss := strings.SplitN(str, " ", -1)

    return strings.Join(ss, "-")
}
