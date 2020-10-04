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
	"fmt"
	"io"
	"io/ioutil"
	heavenshelp "miniWebSrv/Utils" // subdir
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

var appData struct {
	portNumber uint
}

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
		fmt.Fprintf(os.Stderr, "Error reading body: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Println("Body: vvvvvvvvvvvvvvvvvvvv\n", string(body))
	fmt.Println("Body: ^^^^^^^^^^^^^^^^^^^^")
	io.WriteString(w, "Hello, this is miniLogSrv at "+time.Now().UTC().Format(time.RFC3339)+"\r\n")
}

// commandLineOptions just separates the definition of command line options ==> creating a shorter main
func commandLineOptions() []cli.Flag {
	return []cli.Flag{
		cli.UintFlag{
			Name:        "port, p",
			Usage:       "MANDATORY: Port number to listen to. Range: [0..65535]",
			Destination: &appData.portNumber,
		},
	}
}

func main() {
	var err error
	app := cli.NewApp()
	app.Flags = commandLineOptions()
	app.Name = "miniWebSrv"
	app.Version = "0.5.3" // semantic versioning
	app.Usage = "Web Server for testing/echoing the input."
	app.Action = func(c *cli.Context) error {
		heavenshelp.LogInit(app.Name)
		heavenshelp.LogInfo(app.Name + "::version:" + app.Version + ":service starting at: " + time.Now().String())
		if c.Bool("debug") {
			heavenshelp.CondDebugSet(true)
		}
		heavenshelp.CondDebug("Debug is enabled.")
		if appData.portNumber >= 65536 || appData.portNumber == 0 {
			fmt.Fprintln(os.Stderr, "The specified port number must be between 1 and 65535:")
			os.Exit(4)
		}
		heavenshelp.LogInfo("Listening on port: " + fmt.Sprint(appData.portNumber))
		// enable the server
		http.HandleFunc("/", requestHandler)
		err = http.ListenAndServe(":"+strconv.FormatUint(uint64(appData.portNumber), 10), nil)
		return err
	}
	err = app.Run(os.Args)
	if err != nil {
		heavenshelp.LogPanic(err.Error())
	}
}
