package utils

import (
    "time"
    "crypto/rand"
    r "math/rand"
)

func RandomCreateBytes(n int, alphabets ...byte) []byte {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    var randby bool
    if num, err := rand.Read(bytes); num != n || err != nil {
        r.Seed(time.Now().UnixNano())
        randby = true
    }
    for i, b := range bytes {
        if len(alphabets) == 0 {
            if randby {
                bytes[i] = alphanum[r.Intn(len(alphanum))]
            } else {
                bytes[i] = alphanum[b%byte(len(alphanum))]
            }
        } else {
            if randby {
                bytes[i] = alphabets[r.Intn(len(alphabets))]
            } else {
                bytes[i] = alphabets[b%byte(len(alphabets))]
            }
        }
    }
    return bytes
}
