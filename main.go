// Copyright (c) 2020 engel-ch@outlook.com

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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	heavenshelp "miniWebSrv/Utils" // subdir
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var appData struct {
	portNumber uint
	appVersion string
	pathMode   bool
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW REQUEST =======================================================")
	fmt.Println("Method:", r.Method)
	fmt.Println("URL.Path:", r.URL.Path)
	fmt.Println("Host:", r.Host)
	fmt.Println("Header:", r.Header)
	fmt.Println("TrailingHeader:", r.Trailer)
	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("RequestURI:", r.RequestURI)
	fmt.Println("Content-Length:", r.ContentLength)

	if appData.pathMode {
		io.WriteString(w, r.URL.Path+"\r\n")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading body: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body)) // reset reader for ParseForm

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error calling ParseForm...")
	}
	fmt.Println("Form:", r.Form)
	for k := range r.PostForm { // remove post values from Form => Form: getValues, PostForm: PostValues
		r.Form.Del(k)
	}
	fmt.Println("PostForm:", r.PostForm)

	fmt.Println("Body: vvvvvvvvvvvvvvvvvvvv\n", string(body))
	fmt.Println("Body: ^^^^^^^^^^^^^^^^^^^^")
	io.WriteString(w, "Hello, this is miniLogSrv version "+appData.appVersion+" at "+time.Now().UTC().Format(time.RFC3339)+"\r\n")
	io.WriteString(w, "Method: "+r.Method+"\r\n")
	io.WriteString(w, "URL: "+r.URL.Path+"\r\n")
	io.WriteString(w, "Host: "+r.Host+"\r\n")
	io.WriteString(w, "Header:\r\n")
	for key := range r.Header {
		io.WriteString(w, "  "+key+": ")
		for value := range r.Header[key] {
			io.WriteString(w, "    "+r.Header[key][value]+"\r\n")
		}
	}

	io.WriteString(w, "Remote Address: "+r.RemoteAddr+"\r\n")
	io.WriteString(w, "RequestURI: "+r.RequestURI+"\r\n")
	io.WriteString(w, "Content-Length: "+strconv.FormatInt(r.ContentLength, 10)+"\r\n")
	io.WriteString(w, "Form:\r\n")
	for key := range r.Form {
		for value := range r.Form[key] {
			io.WriteString(w, " "+key+": ")
			io.WriteString(w, r.Form[key][value]+"\r\n")
		}
	}
	io.WriteString(w, "PostForm:\r\n")
	for key := range r.PostForm {
		for value := range r.PostForm[key] {
			io.WriteString(w, " "+key+": ")
			io.WriteString(w, r.PostForm[key][value]+"\r\n")
		}
	}
	io.WriteString(w, "Body:\r\n")
	io.WriteString(w, " "+strings.Replace(string(body), "\n", "\n  ", -1)+"\r\n")
}

// commandLineOptions just separates the definition of command line options ==> creating a shorter main
func commandLineOptions() []cli.Flag {
	return []cli.Flag{
		cli.UintFlag{
			Name:        "port, p",
			Usage:       "MANDATORY: Port number to listen to. Range: [0..65535]",
			Destination: &appData.portNumber,
		},
		cli.BoolFlag{
			Name:        "path-mode, m",
			Usage:       "Just return the path. For simple integration into monitor solutions.",
			Destination: &appData.pathMode,
		},
	}
}

func main() {
	var err error
	app := cli.NewApp()
	app.Flags = commandLineOptions()
	app.Name = "miniWebSrv"
	app.Version = "0.12.1" // semantic versioning
	appData.appVersion = app.Version
	app.Usage = "Web Server for testing/echoing the input."
	app.Action = func(c *cli.Context) error {
		fmt.Fprintln(os.Stderr, app.Name+"::version:"+app.Version+":service starting at: "+time.Now().String())
		if c.Bool("debug") {
			heavenshelp.CondDebugSet(true)
		}
		heavenshelp.CondDebug("Debug is enabled.")
		if appData.portNumber >= 65536 || appData.portNumber == 0 {
			fmt.Fprintln(os.Stderr, "The specified port number must be between 1 and 65535:")
			os.Exit(4)
		}
		if appData.pathMode {
			fmt.Fprintln(os.Stderr, "Path mode.")
		}
		fmt.Fprintln(os.Stderr, "Listening on port: "+fmt.Sprint(appData.portNumber))
		// enable the server
		http.HandleFunc("/", requestHandler)
		err = http.ListenAndServe(":"+strconv.FormatUint(uint64(appData.portNumber), 10), nil)
		return err
	}
	err = app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(99)
	}
}
