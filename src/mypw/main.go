// A clone of pwgen.
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func mypw(length int) string {
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
	var length = 20
	var num int = 1
    for i := 0; i < len(os.Args); i++ {
        switch os.Args[i] {
        case "-h", "--help":
            fmt.Printf("Usage: mypw [length=%d] [num=%d]\n", length, num)
            fmt.Println()
		    os.Exit(0)
        }
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
	for i := 0; i < num; i++ {
		fmt.Println(mypw(length))
	}
}
