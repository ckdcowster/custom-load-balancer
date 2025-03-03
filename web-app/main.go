package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Could not get hostname", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello from %s!", hostname)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
