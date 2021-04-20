package fsconvert

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
)

func TestMQTT(t *testing.T) {
	var out io.Writer
	if testing.Verbose() {
		out = os.Stdout
	} else {
		out = io.Discard
	}
	conn, err := net.Dial("tcp", "10.13.13.1:1883")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	c := mqtt.NewClientConn(conn)
	err = c.Connect("", "")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	c.Subscribe([]proto.TopicQos{proto.TopicQos{Topic: "#", Qos: proto.QosAtMostOnce}})
	timeout := time.After(2 * time.Second)
	for {
		select {
		case m := <-c.Incoming:
			fmt.Fprintf(out, "%s = ", m.TopicName)
			m.Payload.WritePayload(out)
			fmt.Fprintln(out)
		case <-timeout:
			return
		}
	}
}
