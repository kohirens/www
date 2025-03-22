package session

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"os"
	"path/filepath"
)

// NewStorageLocal Initialize local session storage.
func NewStorageLocal(workDir string) *LocalStorage {
	return &LocalStorage{
		WorkDir: workDir,
	}
}

type LocalStorage struct {
	WorkDir string
}

// Load Session data from local file storage.
func (ls *LocalStorage) Load(id string) (*Data, error) {
	f := filepath.Join(ls.WorkDir, id)
	if !fsio.Exist(f) {
		return nil, fmt.Errorf("file %v not found", f)
	}

	content, e1 := os.ReadFile(f)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ReadFile, f, e1.Error())
	}

	data := &Data{}

	if e := json.Unmarshal(content, data); e != nil {
		fmt.Printf("JSON to decode: %v\n", string(content))
		return nil, fmt.Errorf(stderr.DecodeJSON, f, e)
	}

	return data, nil
}

// Save Session data to a local file for storage.
func (ls *LocalStorage) Save(data *Data) error {
	f := filepath.Join(ls.WorkDir, data.Id)

	content, e2 := json.Marshal(data)
	if e2 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e2)
	}

	if e := os.WriteFile(f, content, fsio.DefaultFilePerms); e != nil {
		return fmt.Errorf(stderr.WriteFile, e)
	}

	return nil
}
