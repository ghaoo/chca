package utils

func Contains(sl []string, s string) bool {
	contains := false

	for _, v := range sl {
		if v == s {
			contains = true
		}
	}

	return contains
}

func Count(sl []string) (num int) {
	num = 0
	for _, s := range sl {
		if s != "" {
			num += 1
		}
	}

	return
}

func Lt(a, b int) bool { return a < b }
func Eq(a, b int) bool { return a == b }
func Gt(a, b int) bool { return a > b }
