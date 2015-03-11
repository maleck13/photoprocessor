package main

import (
	"fmt"
	"os"
	"testing"
)

func TestReadExif(t *testing.T) {
	m, err := ReadExifData("./testFixtures/DSCF1283.JPG")
	dir, _ := os.Getwd()
	if err != nil {
		t.Errorf(" error getting exif date %s  %s", err.Error(), dir)
	}

	for k, v := range m {
		fmt.Println(" key is : " + k + " val = " + v)
	}
}
