package fsconvert

import (
	"embed"
	"io"
	"os"
	"testing"
)

//go:embed [a-z]*
var content embed.FS

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
