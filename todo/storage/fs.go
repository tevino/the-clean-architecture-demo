package storage

import "github.com/tevino/the-clean-architecture-demo/todo/entity"

// FileSystem implements a file system based storage.
type FileSystem struct {
	path string
}

// NewFileSystem creates a FileSystem with given data path.
func NewFileSystem(path string) *FileSystem {
	return &FileSystem{path}
}

// SaveItem is not implemented yet.
func (f *FileSystem) SaveItem(item *entity.Item) (int64, error) {
	return 0, nil
}

// GetItemsByParentID is not implemented yet.
func (f *FileSystem) GetItemsByParentID(parentID int64) ([]*entity.Item, error) {
	return nil, nil
}

// GetItemByID is not implemented yet.
func (f *FileSystem) GetItemByID(int64) (*entity.Item, error) {
	return nil, nil
}
