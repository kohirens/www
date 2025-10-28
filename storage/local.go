package storage

import (
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"io/fs"
	"os"
)

var ps = string(os.PathSeparator)

// LocalStorage Save data in local files.
type LocalStorage struct {
	WorkDir string
}

func NewLocalStorage(wd string) (*LocalStorage, error) {
	if !fsio.DirExist(wd) {
		return nil, &ErrDirNoExist{wd}
	}

	return &LocalStorage{
		WorkDir: wd,
	}, nil
}

// Load Retrieve file from storage.
func (s *LocalStorage) Load(filename string) ([]byte, error) {
	filePath := s.filePath(filename)

	Log.Dbugf(stdout.Load, filePath)

	if !fsio.Exist(filePath) {
		return nil, fmt.Errorf("%v %v", filePath, fs.ErrNotExist)
	}

	content, e1 := os.ReadFile(filePath)
	if e1 != nil {
		return nil, &ErrReadFile{filePath + " " + e1.Error()}
	}

	return content, nil
}

// Save Write session data to the storage medium.
func (s *LocalStorage) Save(filename string, data []byte) error {
	filePath := s.filePath(filename)

	if e := os.WriteFile(filePath, data, 0774); e != nil {
		return &ErrWriteFile{e.Error()}
	}

	return nil
}

func (s *LocalStorage) filePath(filename string) string {
	return s.WorkDir + ps + filename
}

// Remove Delete a file from storage.
func (s *LocalStorage) Remove(filename string) error {
	fullFilename := s.filePath(filename)

	if e := os.Remove(fullFilename); e != nil {
		return fmt.Errorf(stderr.RemoveFile, e.Error())
	}

	return nil
}
