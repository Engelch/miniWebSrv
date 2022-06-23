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
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	ce "github.com/engelch/debugerrorce/v2"
	"github.com/urfave/cli/v2"
)

const listenPortOption = "port" // option name (use > 1) to specify a listening port

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

// evaluateDebug evaluates command line arguments
func evaluateDebug(c *cli.Context) {
	if c.Bool("debug") {
		ce.CondDebugSet(true)
	}
	ce.CondDebugln("Debug is enabled.")
}

// evaluatePortNumber evaluates command line arguments
func evaluatePortNumber(c *cli.Context, port *uint64) {
	*port = c.Uint64(listenPortOption)
	if *port == 0 || *port >= 65536 {
		ce.ErrorExit(11, "Listening port not in allowed range:"+fmt.Sprintf("%d", *port))
	}
	ce.CondDebugln("Listing port number is: " + fmt.Sprintf("%d", *port))
}

// evaluateArgs does the CLI argument and option validation
func evaluateArgs(c *cli.Context, port *uint64) {
	if c == nil {
		ce.ErrorExit(10, "apd structure is empty.")
	}
	evaluateDebug(c)
	evaluatePortNumber(c, port)
}

// commandLineOptions just separates the definition of command line options ==> creating a shorter main
func commandLineOptions() []cli.Flag {
	return []cli.Flag{
		&cli.UintFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "MANDATORY: Port number to listen to. Range: [1..65535]",
			Value:   0,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"D"},
			Value:   false,
			Usage:   "debug mode",
		},
	}
}

func main() {
	var err error
	var listeningPort uint64
	app := cli.App{}
	app.Flags = commandLineOptions()
	app.Name = "mini-web-svc"
	app.Version = "1.0.0" // semantic versioning
	app.Usage = "Web Server for testing."
	app.Action = func(c *cli.Context) error {
		evaluateArgs(c, &listeningPort) // exits the app in case of error
		fmt.Fprintln(os.Stderr, app.Name+"::version: "+app.Version+"::service starting at "+time.Now().Format(time.RFC3339))
		// enable the server
		http.HandleFunc("/", requestHandler)
		return http.ListenAndServe("127.0.0.1:"+strconv.FormatUint(listeningPort, 10), nil)
	}
	err = app.Run(os.Args)
	if err != nil {
		ce.ErrorExit(99, "Error from app:"+err.Error())
	}
}
