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
	test.ResetDir(tmpDir+"/accounts", os.ModeDir|os.ModePerm)

	os.Exit(m.Run())
}
