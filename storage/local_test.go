package storage

import (
	"bytes"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		workDir  string
		filename string
		wantErr  bool
	}{
		{
			name:     "TestLocalStorage",
			content:  []byte("test1234"),
			workDir:  tmpDir,
			filename: "test1234.txt",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				WorkDir: tt.workDir,
			}

			err := s.Save(tt.filename, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err2 := s.Load(tt.filename)
			if (err2 != nil) != tt.wantErr {
				t.Errorf("LocalStorage.Load(%v) error = %v, wantErr %v", tt.filename, err2, tt.wantErr)
			}

			if !bytes.Equal(content, tt.content) {
				t.Errorf("LocalStorage.Load(%v) = %v, want %v", tt.filename, content, tt.content)
			}

			if err3 := s.Remove(tt.filename); (err3 != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err3, tt.wantErr)
			}

			_, err4 := os.Stat(tmpDir + "/" + tt.filename)
			if os.IsExist(err4) {
				t.Errorf("LocalStorage file %v should have been removed", tt.filename)
			}
		})
	}
}
