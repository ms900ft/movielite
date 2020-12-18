package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

var Movieserver string

func alive(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	io.WriteString(w, `{"alive": "ok"}`)
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

	url := fmt.Sprintf("%s/%s", Movieserver, url.QueryEscape(strings.Replace(r.URL.Path[1:], " ", "_", -1)))
	err := open.RunWith(url, "vlc")
	log.Printf("opening %s", url)
	if err != nil {
		log.Print(err)
		io.WriteString(w, fmt.Sprintf("%q", err.Error))
	} else {
		io.WriteString(w, `{"message": "ok"}`)
	}
}

func main() {
	Movieserver = os.Args[1]
	log.Printf("using movieserver %s", Movieserver)
	http.HandleFunc("/", handler)
	http.HandleFunc("/alive", alive)
	err := http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
