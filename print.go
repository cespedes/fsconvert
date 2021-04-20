package fsconvert

import (
	"fmt"
	"io"
	"io/fs"
)

// PrintTree displays a representation of a filesystem.
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
		fmt.Fprintf(w, "%s %8d %s %s\n", mode.String(), size, modTime.Format("2006-01-02 15:04:05"), path)
		return nil
	})
	return err
}
