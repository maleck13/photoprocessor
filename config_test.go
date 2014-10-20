package main

import "testing"

func TestConfig(t *testing.T) {
	dir := "/Users/craigbrookes/Pictures";
	LoadConfig()

	if x := CONF.GetPhotoDir(); x != dir {
		t.Errorf(" CONF.GetPhotoDir() got %s should be %s ", x, dir)
	}
}
