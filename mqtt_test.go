package fsconvert

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestMQTT(t *testing.T) {
	var out io.Writer
	if testing.Verbose() {
		out = os.Stdout
	} else {
		out = io.Discard
	}
	mqtt := os.Getenv("MQTT")
	if mqtt == "" {
		return
	}
	fsys, err := FromMQTT(mqtt, "#")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	time.Sleep(1 * time.Second)

	err = PrintTree(fsys, out)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
