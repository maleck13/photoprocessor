package processor

import (
	"os"
	"testing"
)

func TestProcessZip(t *testing.T) {
	zipPath := "./testFixtures/Archive.zip"
	dataDir := "./testFixtures/data"
	user := "testuser"

	os.RemoveAll(dataDir + "/" + user)

	loc, err := ProcessZip(dataDir, zipPath, user)

	if nil != err {
		t.Errorf(" error in process zip %s ", err)
	}

	if loc == "" {
		t.Errorf(" error no file path returned ")
	}

	os.RemoveAll(dataDir + "/" + user)

}
