package main

import (
	"testing"
)

func TestAppData(t *testing.T) {
	if appData.portNumber != 0 {
		t.Error()
	}
}

func TestCommandLineOptions(t *testing.T) {
	clio := commandLineOptions()
	if clio == nil {
		t.Error()
	}
}
