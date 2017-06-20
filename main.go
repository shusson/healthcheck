package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"bytes"
	"errors"
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
	r.HandleFunc("/", Index(r, url))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func Index(router *mux.Router, url string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var jsonStr = []byte(`[1,"connect",1,0,{"1":{"str":""},"2":{"str":""},"3":{"str":""}}]`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

		if check(err, w) != nil {
			return
		}
		req.Header.Set("Accept", "application/vnd.apache.thrift.json; charset=utf-8")
		req.Header.Set("Content-Type", "application/vnd.apache.thrift.json; charset=UTF-8")
		req.Header.Set("Content-Length", "88")

		client := &http.Client{}
		resp, err := client.Do(req)
		if check(err, w) != nil {
			return
		}
		defer resp.Body.Close()

		var status error
		if resp.StatusCode != 200 {
			status = errors.New(resp.Status)
		}
		if check(status, w) != nil {
			return
		}
		w.Header().Set("Server", "Mapd healthcheck")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
	return http.HandlerFunc(fn)
}

func check(err error, w http.ResponseWriter) error {
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return err
}
