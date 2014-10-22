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
	UseDefaultLonLat bool
	DefaultLonLat []float64
	DefaultUser string
}

const (
	CONF_ENV_VAR = "PHOTO_PROC_CONF"
)

var (
	CONF *CONFIG
)

func LoadConfig (){
	confPath := os.Getenv(CONF_ENV_VAR)
	file,err := os.Open(confPath)
	if err != nil{
		ErrorLog.Fatal("failed to load config " + confPath + err.Error())
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

func (c * CONFIG)GetUseDefaultLonLat()bool{
	return c.UseDefaultLonLat
}

func (c * CONFIG)GetDefaultLonLat()[]float64{
	return c.DefaultLonLat;
}

func (c * CONFIG)GetDefaultUser()string{
	return c.DefaultUser;
}
