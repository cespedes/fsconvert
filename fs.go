package fsconvert

import (
	"io/fs"
	"testing/fstest"
	"sync"
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

func (fsys FS) Open(name string) (fs.File, error) {
	fsys.mutex.RLock()
	f, err := fsys.mapFS.Open(name)
	fsys.mutex.RUnlock()
	return file{mutex: &fsys.mutex, mapFile: f}, err
}

func (fsys FS) ReadDir(name string) ([]fs.DirEntry, error) {
	fsys.mutex.RLock()
	de, err := fsys.mapFS.ReadDir(name)
	fsys.mutex.RUnlock()
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
	n, err := f.mapFile.Read(b)
	f.mutex.RUnlock()
	return n, err
}

func (f file) Close() error {
	f.mutex.RLock()
	err := f.mapFile.Close()
	f.mutex.RUnlock()
	return err
}

// Remove removes the named file or (empty) directory.
// It always returns nil, even if the name does not exist
// or cannot be removed (because it belongs to a non-empty directory)
func (fsys *FS) Remove(name string) error {
	fsys.mutex.Lock()
	delete(fsys.mapFS, name)
	fsys.mutex.Unlock()
	return nil
}

// WriteFile creates or overwrites the named file with the specified
// data and permissions.
func (fsys *FS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	mapFile := fstest.MapFile{Data: data, Mode: perm, ModTime: time.Now()}
	fsys.mutex.Lock()
	if fsys.mapFS == nil {
		fsys.mapFS = make(fstest.MapFS)
	}
	fsys.mapFS[name] = &mapFile
	fsys.mutex.Unlock()
	return nil
}
