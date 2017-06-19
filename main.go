package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var url string
	flag.StringVar(&url, "u", "http://localhost:80", "url to mapd-core server")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if url == "" {
		flag.PrintDefaults()
		return
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", Index(r))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func Index(router *mux.Router) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Mapd healthcheck")
	}
	return http.HandlerFunc(fn)
}

func check(err error) {
	if err == nil {
		return
	}
	panic(err)
}

