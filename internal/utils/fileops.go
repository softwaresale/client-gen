package utils

import "os"

type CompilerFileOperations interface {
	Create(filename string) (*os.File, error)
	Stat(filename string) (os.FileInfo, error)
	Mkdir(filename string, perm os.FileMode) error
	IsNotExist(err error) bool
}

// OSFileOperations provide the standard file operations interface
type OSFileOperations struct{}

func (O OSFileOperations) Create(filename string) (*os.File, error) {
	return os.Create(filename)
}

func (O OSFileOperations) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

func (O OSFileOperations) Mkdir(filename string, perm os.FileMode) error {
	return os.Mkdir(filename, perm)
}

func (O OSFileOperations) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
