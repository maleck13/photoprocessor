package main

import (
	"encoding/json"
	"os"
	"strings"
	"io/ioutil"
	"fmt"
)

type CONFIG struct {
	PhotoDir , MongoHost, ProcessedPhotoDir, ThumbnailDir string
	ConcurrentJobs int
}

const (
	CONF_PATH = "/etc/photoprocessor/conf.json"
)

var (
	CONF *CONFIG
)

func LoadConfig (){
	file,err := os.Open(CONF_PATH)
	if err != nil{
		ErrorLog.Fatal("failed to load config " + CONF_PATH + err.Error())
	}
	contentBuf,err := ioutil.ReadAll(file)
	dec := json.NewDecoder(strings.NewReader(string(contentBuf)))
	var c CONFIG
	CONF =&c;

	err = dec.Decode(CONF)
	FailOnError(err, "failed to decode config")
	fmt.Println(CONF);

}


func (c * CONFIG)GetPhotoDir()string{
	return c.PhotoDir
}

func (c * CONFIG)GetMongoHost()string{
	return c.MongoHost
}

func (c * CONFIG)GetConcurrentJobs()int{
	return c.ConcurrentJobs
}

func (c * CONFIG)GetProcessedPhotoDir()string{
	return c.ProcessedPhotoDir
}

func (c * CONFIG)GetThumbNailDir()string{
	return c.ThumbnailDir;
}
