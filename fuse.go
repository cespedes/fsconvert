package fsconvert

import (
	"context"
	"hash/crc64"
	"io/fs"
	"log"
	"os"
	"path"
	"syscall"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
)

// ToFUSE mountes a filesystem using FUSE on the named mountpoint.
func ToFUSE(fsys fs.FS, mountpoint string) error {
	log.Printf("ToFUSE: Will mount filesystem to %s\n", mountpoint)
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

func (f fuseFS) Root() (fusefs.Node, error) {
	return fuseDir{fsys: f.fsys, path: "."}, nil
}

type fuseDir struct {
	fsys fs.FS
	path string
}

func fuseCalcInode(name string) uint64 {
	return crc64.Checksum([]byte(name), crc64.MakeTable(crc64.ISO))
}

func (dir fuseDir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = fuseCalcInode(dir.path)
	log.Printf("FUSE:Attr(dir=%s): inode=%d\n", dir.path, a.Inode)
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (dir fuseDir) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	log.Printf("FUSE:Lookup(dir=%s, name=%s)\n", dir.path, name)
	newpath := path.Join(dir.path, name)
	fi, err := fs.Stat(dir.fsys, newpath)
	if err != nil {
		return nil, syscall.ENOENT
	}
	if fi.IsDir() {
		return fuseDir{fsys: dir.fsys, path: newpath}, nil
	}
	return fuseFile{fsys: dir.fsys, path: newpath}, nil
}

func (dir fuseDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Printf("FUSE:ReadDirAll(dir=%s)\n", dir.path)
	entries, err := fs.ReadDir(dir.fsys, dir.path)
	if err != nil {
		return nil, err
	}
	var result []fuse.Dirent
	for _, e := range entries {
		name := e.Name()
		isDir := e.IsDir()
		var t fuse.DirentType = fuse.DT_File
		if isDir {
			t = fuse.DT_Dir
		}
		result = append(result, fuse.Dirent{Inode: fuseCalcInode(path.Join(dir.path, name)), Name: name, Type: t})
	}
	return result, nil
}

type fuseFile struct {
	fsys fs.FS
	path string
}

func (file fuseFile) Attr(ctx context.Context, a *fuse.Attr) error {
	fi, err := fs.Stat(file.fsys, file.path)
	if err != nil {
		return syscall.ENOENT
	}
	a.Inode = fuseCalcInode(file.path)
	a.Mode = fi.Mode()
	a.Size = uint64(fi.Size())
	return nil
}

func (file fuseFile) ReadAll(ctx context.Context) ([]byte, error) {
	return fs.ReadFile(file.fsys, file.path)
}

/*

// fuseFile implements both Node and Handle for the hello file.

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
