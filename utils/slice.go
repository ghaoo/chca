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
