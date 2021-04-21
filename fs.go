package fsconvert

import (
	"io/fs"
	"sync"
	"testing/fstest"
	"time"
)

// FS is a simple structure implementing io/fs's FS interface
// It is implemented as a MapFS (from testing/fstest) but with locks
// to be able to have concurrent reads and writes
type FS struct {
	mutex sync.RWMutex
	mapFS fstest.MapFS
}

type file struct {
	mutex   *sync.RWMutex
	mapFile fs.File
}

// Open opens the named file.
func (fsys FS) Open(name string) (fs.File, error) {
	fsys.mutex.RLock()
	defer fsys.mutex.RUnlock()
	f, err := fsys.mapFS.Open(name)
	return file{mutex: &fsys.mutex, mapFile: f}, err
}

// ReadDir opens the named directory and returns a list of
// directory entries sorted by filename.
func (fsys FS) ReadDir(name string) ([]fs.DirEntry, error) {
	fsys.mutex.RLock()
	defer fsys.mutex.RUnlock()
	de, err := fsys.mapFS.ReadDir(name)
	return de, err
}

func (f file) Stat() (fs.FileInfo, error) {
	f.mutex.RLock()
	fi, err := f.mapFile.Stat()
	f.mutex.RUnlock()
	return fi, err
}

func (f file) Read(b []byte) (int, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	n, err := f.mapFile.Read(b)
	return n, err
}

func (f file) Close() error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	err := f.mapFile.Close()
	return err
}

// Remove removes the named file or (empty) directory.
// It may return nil, even if the name does not exist
// or cannot be removed (because it belongs to a non-empty directory)
func (fsys *FS) Remove(name string) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrNotExist}
	}
	fsys.mutex.Lock()
	defer fsys.mutex.Unlock()
	delete(fsys.mapFS, name)
	return nil
}

// WriteFile creates or overwrites the named file with the specified
// data and permissions.
func (fsys *FS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "WriteFile", Path: name, Err: fs.ErrNotExist}
	}
	mapFile := fstest.MapFile{Data: data, Mode: perm, ModTime: time.Now()}
	fsys.mutex.Lock()
	defer fsys.mutex.Unlock()
	if fsys.mapFS == nil {
		fsys.mapFS = make(fstest.MapFS)
	}
	fsys.mapFS[name] = &mapFile
	return nil
}

// AppendFile appends data to the already existing named file.
func (fsys *FS) AppendFile(name string, data []byte) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "AppendFile", Path: name, Err: fs.ErrNotExist}
	}
	fsys.mutex.Lock()
	defer fsys.mutex.Unlock()
	if fsys.mapFS == nil || fsys.mapFS[name] == nil {
		return &fs.PathError{Op: "AppendFile", Path: name, Err: fs.ErrNotExist}
	}
	mapFile := fsys.mapFS[name]
	mapFile.Data = append(mapFile.Data, data...)
	mapFile.ModTime = time.Now()
	return nil
}
