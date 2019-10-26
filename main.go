// Copyright (c) 2019 engel-ch@outlook.com

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// ================================================================
// miniWebSrv :- simplistic web-server for debugging purposes.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"io/ioutil"
)

const appVersion = "0.4.1"		// semantic versioning

func requestHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintln(os.Stderr, "Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Println("Body: vvvvvvvvvvvvvvvvvvvv\n", string(body))
	fmt.Println("Body: ^^^^^^^^^^^^^^^^^^^^")
	io.WriteString(w, "Hello, this is miniLogSrv at "+time.Now().UTC().Format(time.RFC3339)+"\r\n")
}

func main() {
	// cmdLine parsing
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
	// enable the server
	http.HandleFunc("/", requestHandler)
	err = http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:ListenAndServe:" + err.Error())
		os.Exit(9)
	}
}
