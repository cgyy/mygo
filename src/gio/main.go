package main

import (
	"encoding/json"
	"flag"
	"gio/functions"
	"gio/routes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

func funcwrapper(myfunc func(input interface{}) interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		var f interface{}
		err = json.Unmarshal(data, &f)
		if err != nil {
			panic(err)
		}
		m := myfunc(f)
		b, err := json.Marshal(m)
		w.Write(b)
	}
}

func main() {
	daemon := flag.Bool("daemon", false, "Indicate it's daemon process. Never use it in command line.")
	flag.Parse()

	if *daemon {
		fd, err := syscall.Open("/dev/null", syscall.O_RDWR, 0)

		if err != nil {
			panic(err)
		}

		syscall.Dup2(fd, syscall.Stdin)
		syscall.Dup2(fd, syscall.Stdout)
		syscall.Dup2(fd, syscall.Stderr)

		if fd > syscall.Stderr {
			syscall.Close(fd)
		}
	}

	// start gio as daemon process
	if !*daemon {
		args := append([]string{os.Args[0], "-daemon"}, os.Args[1:]...)
		attr := syscall.ProcAttr{}
		_, _, err := syscall.StartProcess(os.Args[0], args, &attr)

		if err != nil {
			panic(err)
			return
		}
		return
	}

	router := routes.NewRouter()
	for _, function := range functions.Functions {
		router.Get(function.Path, funcwrapper(function.Func))
		router.Post(function.Path, funcwrapper(function.Func))
	}
	s := &http.Server{
		Addr:           ":1234",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
