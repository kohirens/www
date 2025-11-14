package backend

import (
	"github.com/kohirens/stdlib/test"
	"os"
	"testing"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
)

func TestMain(m *testing.M) {
	test.ResetDir(tmpDir, 0777)
	test.ResetDir(tmpDir+"/accounts", 0777)

	os.Exit(m.Run())
}
