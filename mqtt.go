package fsconvert

import (
	"io/fs"
	_ "github.com/jeffallen/mqtt"
)

func FromMQTT(address string, topic string) fs.FS {
	return FS{}
}
