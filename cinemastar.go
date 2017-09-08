package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var base string
var port string

func main() {
	flag.StringVar(&base, "f", "", "SearchPattern")
	flag.StringVar(&port, "p", "8080", "port")
	flag.Parse()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/latest", Latest)
	router.GET("/file/*filepath", MovieIndex)
	router.GET("/category/*filepath", MovieCategoryIndex)
	router.GET("/static/*filepath", MovieContent)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
