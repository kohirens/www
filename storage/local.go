package storage

import (
	"github.com/kohirens/stdlib/fsio"
	"io/fs"
	"os"
)

var ps = string(os.PathSeparator)

// LocalStorage Save data in local files.
type LocalStorage struct {
	Name    string
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

// Load data from a file in to the data store.
func (s *LocalStorage) Load(filename string) ([]byte, error) {
	filePath := s.WorkDir + ps + filename

	Log.Dbugf("load %v", filePath)

	if !fsio.Exist(filePath) {
		return nil, fs.ErrNotExist
	}

	content, e1 := os.ReadFile(filePath)
	if e1 != nil {
		return nil, &ErrReadFile{filePath + " " + e1.Error()}
	}

	return content, nil
}

// Save The session data to the storage medium.
func (s *LocalStorage) Save(filename string, data []byte) error {
	filePath := s.WorkDir + ps + filename

	if e := os.WriteFile(filePath, data, 0744); e != nil {
		return &ErrWriteFile{e.Error()}
	}

	return nil
}
