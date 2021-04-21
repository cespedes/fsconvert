package main

import (
	"fmt"
	"os"

	"github.com/cespedes/fsconvert"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: mqtt2fuse <mqtt-server> <topic> <mountpoint>")
		os.Exit(1)
	}
	mqtt := os.Args[1]
	topic := os.Args[2]
	mountpoint := os.Args[3]

	fsys, err := fsconvert.FromMQTT(mqtt, topic)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = fsconvert.ToFUSE(fsys, mountpoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
