package main

import (
	"fmt"
	"log"
	"net/http"
	"zweb"
)

func main() {
	r := zweb.New()
	r.GET("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})
	r.GET("/hello", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			_, _ = fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
		}
	})

	log.Fatal(r.Run(":10002"))
}
