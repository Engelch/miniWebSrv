package heavenshelp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
)

// threat-safe implementation
// todo hide variable by closure, PRIO: low
// globalDebug is used to store the state if debugging is on/off. This variable should
// never be set directly.
var globalDebug atomic.Value

// init() is always executed at the startup of the application. It makes sure that
// the global debug functionality has a defined state: off ⇔ false.
func init() {
	globalDebug.Store(false)
}

// CondDebug() is the implementation of a global debug function. If it was turned on using
// CondDebugSet(true), then the string is shown to stderr. Else, no output is created.
func CondDebug(msg string) {
	if globalDebug.Load().(bool) {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// CondDebugSet(val bool) allows us to turn debug on/off.
func CondDebugSet(val bool) {
	globalDebug.Store(val)
}

// CondDebugStatus() allows to check if debug is turned on/off.
func CondDebugStatus() bool {
	return globalDebug.Load().(bool)
}

// CaptureOutput get a function as its argument. It executes the function and returns the output (stderr and stdout) created by this function.
// While capturing this output, this output is not written to default stdout or stderr.
func CaptureOutput(f func()) (stderr string, stdout string) {
	//fmt.Println("in captureOutput")
	rerr, werr, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating error pipe\n")
		os.Exit(1)
	}
	rout, wout, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output pipe\n")
		os.Exit(1)
	}

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	olderr := os.Stderr
	oldout := os.Stdout

	os.Stderr = werr
	os.Stdout = wout
	f()
	werr.Close()
	wout.Close()

	os.Stderr = olderr
	os.Stdout = oldout
	io.Copy(errbuf, rerr)
	io.Copy(outbuf, rout)

	//fmt.Printf("Copied bytes: %d\n", buf.Len())
	//fmt.Printf("Contents: %s\n", string(buf.Bytes()))
	rerr.Close()
	rout.Close()
	return string(errbuf.Bytes()), string(outbuf.Bytes())
}

// CurrentFunctionName() returns the name of the current function being executed.
func CurrentFunctionName() string {
	pc := make([]uintptr, 1) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
