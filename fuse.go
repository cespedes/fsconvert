package fsconvert

import (
	"context"
	"io/fs"
	"log"
	"os"
	//"syscall"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
)

// ToFUSE mountes a filesystem using FUSE on the named mountpoint.
func ToFUSE(fsys fs.FS, mountpoint string) error {
	log.Printf("Will mount filesystem to %s\n", mountpoint)
	c, err := fuse.Mount(mountpoint, fuse.FSName("fsconvert"), fuse.Subtype("ToFUSE"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fusefs.Serve(c, fuseFS{fsys: fsys})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// fuseFS implements the hello world file system.
type fuseFS struct {
	fsys fs.FS
}

func (fuseFS) Root() (fusefs.Node, error) {
	return fuseDir{}, nil
}

// fuseDir implements both Node and Handle for the root directory.
type fuseDir struct{}

func (fuseDir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}

/*
func (fuseDir) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	if name == "hello" {
		return fuseFile{}, nil
	}
	return nil, syscall.ENOENT
}

func (fuseDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var dirDirs = []fuse.Dirent{
		{Inode: 2, Name: "hello", Type: fuse.DT_File},
	}
	return dirDirs, nil
}

// fuseFile implements both Node and Handle for the hello file.
type fuseFile struct{}

const fuseGreeting = "hello, world\n"

func (fuseFile) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0o444
	a.Size = uint64(len(fuseGreeting))
	return nil
}

func (fuseFile) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(fuseGreeting), nil
}
*/
