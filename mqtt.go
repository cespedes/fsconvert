package fsconvert

import (
	"bytes"
	"io/fs"
	"net"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
)

// FromMQTT creates a filesystem, filling it with the values published
// to a MQTT server under a specified topic.
func FromMQTT(address string, topic string) (fs.FS, error) {
	var fsys FS

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	c := mqtt.NewClientConn(conn)
	err = c.Connect("", "")
	if err != nil {
		return nil, err
	}

	c.Subscribe([]proto.TopicQos{proto.TopicQos{Topic: "#", Qos: proto.QosAtMostOnce}})

	/* {
		var b bytes.Buffer
		m := <-c.Incoming
		m.Payload.WritePayload(&b)
		fsys.WriteFile(m.TopicName, b.Bytes(), 0444)
	} */

	go func() {
		for {
			var b bytes.Buffer
			m := <-c.Incoming
			m.Payload.WritePayload(&b)
			fsys.WriteFile(m.TopicName, b.Bytes(), 0444)
		}
	}()
	return &fsys, nil
}
