package log

import (
	"io"
	"os"
)

type FileSystemInterface interface {
	Create(fileName string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
	Remove(name string) error
	WriteString(w io.Writer, s string) (n int, err error)
}

////////////////////////////////////////////////////////////////////////////////
// Real structure

type FileSystem struct {
}

func (f FileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
func (f FileSystem) Create(filename string) (*os.File, error)             { return os.Create(filename) }
func (f FileSystem) Stat(name string) (os.FileInfo, error)                { return os.Stat(name) }
func (f FileSystem) Remove(name string) error                             { return os.Remove(name) }
func (f FileSystem) WriteString(w io.Writer, s string) (n int, err error) { return io.WriteString(w, s) }

////////////////////////////////////////////////////////////////////////////////
// Mock structure

type MockFileSystem struct {
	mockStat        error
	mockCreate      error
	mockRemove      error
	mockWriteString error
	mockOpenFile    error
}

func (f MockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, f.mockOpenFile
}
func (f MockFileSystem) Create(filename string) (*os.File, error) { return nil, f.mockCreate }
func (f MockFileSystem) Stat(name string) (os.FileInfo, error)    { return nil, f.mockStat }
func (f MockFileSystem) Remove(name string) error                 { return f.mockRemove }
func (f MockFileSystem) WriteString(w io.Writer, s string) (n int, err error) {
	return 0, f.mockWriteString
}
