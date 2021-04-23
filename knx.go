package fsconvert

import (
	"fmt"
	"io/fs"

	"github.com/vapourismo/knx-go/knx"
)

// FromKNX creates a filesystem, filling it with the values sent
// to KNX Group Addresses.  Source of the packets is not stored.
func FromKNX(router string) (fs.FS, error) {
	client, err := knx.NewGroupTunnel(router, knx.DefaultTunnelConfig)
	if err != nil {
		return nil, err
	}

	knxChan := client.Inbound()

	var fsys FS

	go func() {
		defer client.Close()
		for {
			event := <-knxChan
			fmt.Printf("%s = %v\n", event.Destination.String(), event.Data)
			fsys.WriteFile(event.Destination.String(), event.Data, 0444)
		}
	}()
	return &fsys, nil
}
