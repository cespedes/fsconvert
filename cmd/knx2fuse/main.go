package main

import (
	"fmt"
	"os"

	"github.com/cespedes/fsconvert"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: knx2fuse <knx-server> <mountpoint>")
		os.Exit(1)
	}
	knx := os.Args[1]
	mountpoint := os.Args[2]

	fsys, err := fsconvert.FromKNX(knx)
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
