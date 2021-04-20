package fsconvert

// FS is a simple structure implementing io/fs's FS interface
type FS struct {
}

type FSchange int

const (
	FSTypeAdd = iota
	FSTypeDel
	FSTypeChange
)
