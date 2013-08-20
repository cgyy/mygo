/**
 * http://golang.org/doc/effective_go.html#web_server
 * with groupcache (https://github.com/golang/groupcache)
 */

package main

import (
	"flag"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	gc "github.com/golang/groupcache"
	"encoding/base64"
)

var addr = flag.String("addr", ":1718", "http service address")
var port = flag.String("port", "8001", "groupcache http port")

var templ = template.Must(template.New("qr").Parse(templateStr))
var cg *gc.Group

func main() {
	flag.Parse()

	peers := gc.NewHTTPPool("http://localhost:" + *port)
	peers.Set("http://localhost:8001", "http://localhost:8002")

	cg = gc.NewGroup("QRCache", 1<<20, gc.GetterFunc(
		func(ctx gc.Context, key string, dest gc.Sink) error {
			fmt.Printf("asking for data of %s\n", key)
			url := fmt.Sprintf("http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl=%s", key)
			resp, err := http.Get(url)
			if err != nil {
				return nil
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil
			}
			value := base64.StdEncoding.EncodeToString(body)
			dest.SetBytes([]byte(value))
			return nil
		}))

	go http.ListenAndServe("localhost:" + *port, peers)

	http.Handle("/", http.HandlerFunc(QR))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func QR(w http.ResponseWriter, req *http.Request) {
	locals := make(map[string]string)
	locals["s"] = req.FormValue("s")
	if len(locals["s"]) > 0 {
		var value string
		if err := cg.Get(nil, locals["s"], gc.StringSink(&value)); err != nil {
			fmt.Printf("err: %v\n", err)
		}
		locals["image"] = value
	}
	templ.Execute(w, locals)
	// show groupcache stats
	fmt.Printf("####### Stats ######")
	fmt.Printf("Group Stats:\n")
	fmt.Printf("   Gets: %d\n", cg.Stats.Gets)
	fmt.Printf("   CacheHits: %d\n", cg.Stats.CacheHits)
	fmt.Printf("   PeerLoads: %d\n", cg.Stats.PeerLoads)
	fmt.Printf("   PeerErrors: %d\n", cg.Stats.PeerErrors)
	fmt.Printf("   Loads: %d\n", cg.Stats.Loads)
	fmt.Printf("   LoadsDeduped: %d\n", cg.Stats.LoadsDeduped)
	fmt.Printf("   LocalLoads: %d\n", cg.Stats.LocalLoads)
	fmt.Printf("   LocalLoadErrs: %d\n", cg.Stats.LocalLoadErrs)
	fmt.Printf("   ServerRequests: %d\n", cg.Stats.ServerRequests)
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .s}}
<img src="data:image/png;base64,{{.image}}" />
<br>
{{.s}}
<br>
<br>
{{end}}
<form action="/" method="GET">
<input maxLength=1024 size=70 name="s" value="" title="Text to QR Encode">
<input type="submit" value="Show QR">
</form>
</body>
</html>
`
