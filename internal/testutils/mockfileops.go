package testutils

import (
	"io/fs"
	"os"
	"time"
)

// MockFileOperations provides
type MockFileOperations struct {
	TestCreate     func(filename string) (*os.File, error)
	TestStat       func(filename string) (os.FileInfo, error)
	TestMkdir      func(filename string, perm os.FileMode) error
	TestIsNotExist func(err error) bool
}

func DefaultMockFileOperations() *MockFileOperations {
	return &MockFileOperations{
		TestCreate: func(filename string) (*os.File, error) {
			return nil, nil
		},
		TestStat: func(filename string) (os.FileInfo, error) {
			return nil, nil
		},
		TestMkdir: func(filename string, perm os.FileMode) error {
			return nil
		},
		TestIsNotExist: func(err error) bool {
			return false
		},
	}
}

func (m *MockFileOperations) Create(filename string) (*os.File, error) {
	return m.TestCreate(filename)
}

func (m *MockFileOperations) Stat(filename string) (os.FileInfo, error) {
	return m.TestStat(filename)
}

func (m *MockFileOperations) Mkdir(filename string, perm os.FileMode) error {
	return m.TestMkdir(filename, perm)
}

func (m *MockFileOperations) IsNotExist(err error) bool {
	return m.TestIsNotExist(err)
}

type MockFileInfo struct {
	NameV    string
	SizeV    int64
	ModeV    fs.FileMode
	ModTimeV time.Time
	IsDirV   bool
	SysV     any
}

func (m MockFileInfo) Name() string {
	return m.NameV
}

func (m MockFileInfo) Size() int64 {
	return m.SizeV
}

func (m MockFileInfo) Mode() fs.FileMode {
	return m.ModeV
}

func (m MockFileInfo) ModTime() time.Time {
	return m.ModTimeV
}

func (m MockFileInfo) IsDir() bool {
	return m.IsDirV
}

func (m MockFileInfo) Sys() any {
	return m.SysV
}
