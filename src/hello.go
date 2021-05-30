package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":80", nil))

	Put("test", "TEST VALUE")
	val, error := Get("test1")
	if error != nil {
		fmt.Println(error)
		return
	}
	fmt.Println(val)
}
