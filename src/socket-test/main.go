package main

import (
	"net"
    "fmt"
	"log"
    "flag"
)

func main() {
	var port = flag.Int("port", 7919, "")
	var host = flag.String("host", "localhost", "")
	flag.Parse()

    conns := []net.Conn{}
    for i :=0; true; i++ {
        conn, e := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
        if e != nil {
            fmt.Printf("\n")
            log.Fatal(e)
            break
        }
        if i % 100 == 0 {
            fmt.Printf("\r%d", i)
        }
        conns = append(conns, conn)
    }
}
