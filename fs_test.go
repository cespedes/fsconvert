package fsconvert

import (
	"io"
	"os"
	//"io/fs"

	"testing"
)

func TestFS(t *testing.T) {
	var out io.Writer
	if testing.Verbose() {
		out = os.Stdout
	} else {
		out = io.Discard
	}

	var fsys FS
	fsys.WriteFile("uno/one", []byte("1"), 0666)
	fsys.WriteFile("dos/two", []byte("2"), 0444)
	fsys.WriteFile("tres/three", []byte("3"), 0751)
	fsys.WriteFile("foo/bar", []byte{1, 2, 3}, 0)
	fsys.WriteFile("hola", []byte("123"), 0640)

	err := PrintTree(fsys, out)
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
