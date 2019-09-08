package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type myHttpHandler struct{}

func (m myHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, this is miniSrv at "+time.Now().UTC().Format(time.RFC3339)+"\r\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Port # to listen to must be specified as first argument")
		os.Exit(1)
	}
	portStr := os.Args[1]
	portInt, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Argument could not be converted to int:", portStr)
		os.Exit(2)
	}
	if portInt < 0 {
		fmt.Fprintln(os.Stderr, "Argument must be ≥ 0:", portStr)
		os.Exit(4)
	}
	// http.HandleFunc("/", hello)
	var h2 myHttpHandler
	http.ListenAndServe(":"+portStr, h2)
}
