package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"io/ioutil"
	"log"
)

const appVersion = "0.3.1"

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL.Path:", r.URL.Path)
	fmt.Println("Host:", r.Host)
	fmt.Println("Header:", r.Header)
	fmt.Println("TrailingHeader:", r.Trailer)
	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("RequestURI:", r.RequestURI)
	fmt.Println("Content-Length:", r.ContentLength)
	fmt.Println("Form:", r.Form)
	body, err := ioutil.ReadAll(r.Body)
   if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		// todo add http error message
		return
	}
	fmt.Println("Body: vvvvvvvvvvvvvvvvvvvv\n", string(body))
	fmt.Println("Body: ^^^^^^^^^^^^^^^^^^^^")
	io.WriteString(w, "Hello, this is miniLogSrv at "+time.Now().UTC().Format(time.RFC3339)+"\r\n")
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
		os.Exit(3)
	}
	if portInt >= 65536 {
		fmt.Fprintln(os.Stderr, "Argument must be ≤ 65535:", portStr)
		os.Exit(4)
	}
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+portStr, nil)
}
