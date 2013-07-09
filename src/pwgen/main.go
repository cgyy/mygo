package main

import (
    "fmt"
    "os"
    "crypto/rand"
    "math/big"
    "strconv"
)

func pwgen(length int) string {
    if length < 1 {
        length = 10
    }
    buf := make([]byte, length)
    for i := 0; i < length; i++ {
        b, _ := rand.Int(rand.Reader, big.NewInt(3))
        var base byte
        var a *big.Int
        switch b.Int64() {
        // lowercase
        case 0:
            base = 'a'
            a, _ = rand.Int(rand.Reader, big.NewInt(26))
        // uppercase
        case 1:
            base = 'A'
            a, _ = rand.Int(rand.Reader, big.NewInt(26))
        // digit
        case 2:
            base = '0'
            a, _ = rand.Int(rand.Reader, big.NewInt(10))
        }
        buf[i] = base + byte(a.Int64())
    }
    return string(buf[:])
}

func main() {
    var num int = 1
    var length int = 10
    if len(os.Args) <= 1 {
        fmt.Println("Usage: pwgen [length=10] [num=1]")
        fmt.Println()
        os.Exit(0)
    }
    if len(os.Args) > 1 {
        _num, err := strconv.Atoi(os.Args[1])
        if err == nil {
            length = _num
        }
    }
    if len(os.Args) > 2 {
        _num, err := strconv.Atoi(os.Args[2])
        if err == nil {
            num = _num
        }
    }
    for i:= 0; i < num; i++ {
        fmt.Println(pwgen(length))
    }
}
