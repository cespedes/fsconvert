package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
)

//go:embed etc/*
var content embed.FS

func PrintTree(f fs.FS, w io.Writer) error {
	err := fs.WalkDir(f, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// name := de.Name()
		// idDir := de.IsDir()
		// typ := de.Type()
		info, err := de.Info()
		if err != nil {
			return err
		}
		size := info.Size()
		mode := info.Mode()
		modTime := info.ModTime()
		fmt.Fprintf(w, "%s %8d %s %s\n", mode.String(), size, modTime.Format("2006-01-02 15:04"), path)
		return nil
	})
	return err
}

func main() {
	PrintTree(content, os.Stdout)
}
