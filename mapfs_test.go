package fsconvert

import (
	"io"
	"os"
	"time"

	"testing"
	"testing/fstest"
)

func TestMapFS(t *testing.T) {
	var out io.Writer
	if testing.Verbose() {
		out = os.Stdout
	} else {
		out = io.Discard
	}

	mapFS := make(fstest.MapFS)
	mapFS["foo/bar"] = &fstest.MapFile{ModTime: time.Now()}

	err := PrintTree(mapFS, out)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

/*
func TestPrintTree(t *testing.T) {
	var out io.Writer
	if testing.Verbose() {
		out = os.Stdout
	} else {
		out = io.Discard
	}
	err := PrintTree(content, out)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
*/
