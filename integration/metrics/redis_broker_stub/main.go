package main

import (
	"fmt"
	"net/http"
)

func debugHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"pool":{"count":1,"clusters":[["10.10.32.9"]]},"allocated":{"count":0,"clusters":null}}`)
}

func main() {
	http.HandleFunc("/debug/", debugHandler)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
}
