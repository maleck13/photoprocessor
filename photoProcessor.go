package main

import (
	"fmt"
	"github.com/gosexy/exif"
	"io/ioutil"
	"strings"
	"os"
	"strconv"
	"errors"
	"image/jpeg"
	"github.com/nfnt/resize"
	"time"
	"io"
)

const (

	EAST_OR_WEST_LON = "East or West Longitude"
	NORTH_OR_SOUTH_LAT = "North or South Latitude"
	DATE_TIME_KEY = "Date and Time (Original)"
	LONGITUDE = "Longitude"
	LATITUDE = "Latitude"
)




type Worker func(chan int)

func ProcessPhotoDir() {

	fInfo, err := ioutil.ReadDir(CONF.GetPhotoDir())
	FailOnError(err," Error reading dir:" + err.Error())

   // takes filinfo and return a func to be executed uses chan to increment a count
   	buildWorker := func (f os.FileInfo ) func(chan int){
	 	return func (c chan int){

	   		if !f.IsDir() && strings.Contains(f.Name(), ".JPG") {
		   		ProcessImg(f.Name())
	   		}
	   		c <- 1
	 	}
   	}

	//build a slice of Worker jobs to be executed
	jobs:=make([]Worker, len(fInfo))

	for idx, f := range fInfo {
		jobs[idx] = buildWorker(f)
	}

	//execute the jobs ten at a time in a go routine
	go executor(jobs,CONF.GetConcurrentJobs())

}

func ProcessImg(filePath string ) {
	reader := exif.New()
	path := CONF.GetPhotoDir() + "/" + filePath
	completedPath := CONF.GetProcessedPhotoDir() + "/"+filePath
	err := reader.Open(path)

	LogOnError(err, "failed to open " + path)


	tags := reader.Tags
	err = validateLonLat(tags);
	if err !=nil{
		LogOnError(err, "missing data")
		return
	}
	err = validateTime(tags);
	if  err !=nil{
		LogOnError(err, "missing data")
		return
	}


	lonLat := convertDegToDec(tags[LATITUDE],tags[NORTH_OR_SOUTH_LAT],tags[LONGITUDE],tags[EAST_OR_WEST_LON])
	pic := Picture{}
	pic.LonLat = lonLat
	pic.Name = filePath
	pic.Path = path
	thumb,err :=createThumb(path,filePath)
	date :=parseDate(tags[DATE_TIME_KEY])

	if err != nil{
		LogOnError(err, "failed to create thumb ignoring img " + filePath)
		return
	}
	pic.Thumb = thumb
	pic.Time = date
	pic.TimeStamp = date.Unix()
	err=SavePic(pic)
	if err != nil{
		LogOnError(err, "failed to save picture")
	}

	//move to completed
	f,err := os.Open(path)

	FailOnError(err,"could not complete by copying to new loc " + completedPath);
	defer  f.Close()
	fc,err := os.Create(completedPath)

	FailOnError(err,"could not complete by copying to new loc " + completedPath);
	defer  fc.Close()

	_,err = io.Copy(fc,f)

	FailOnError(err, "failed to copy file " + path )




}

func validateLonLat (info map[string]string) error{
	_,ok := info[LONGITUDE]
	if !ok{
		return errors.New("no " + LONGITUDE + " field")
	}
	_,ok = info[LATITUDE]

	if !ok{
		return errors.New("no " + LATITUDE + " field")
	}
	return nil
}

func validateTime(info map[string]string) error{
	_,ok := info[DATE_TIME_KEY]
	if ! ok{
		return errors.New("no " + DATE_TIME_KEY + " present")
	}
	return nil
}

func convertDegToDec(latDeg string, latFlag string, lonDeg string, lonFlag string) []float64{
	//12° 34" 56' = 12 + (34/60) + (56/3600) = 12.582222222222222222222
	bits :=strings.Split(latDeg,",")

	fmt.Print(bits[0] + " : " + bits[1] + " : " + bits[2])
	val1,_ := strconv.ParseFloat(strings.TrimSpace(bits[0]),64)
	val2,_ := strconv.ParseFloat(strings.TrimSpace(bits[1]),64)
	val3,_ := strconv.ParseFloat(strings.TrimSpace(bits[2]),64)
	retFloat := make([]float64,2)
	fmt.Printf(" float vals %f %f %f ", val1, val2,val3)
	latDec:= val1 + (val2/60) + (val3/3600)

	if "N" != latFlag && "E" != latFlag{
		latDec = latDec * -1
	}

	bits = strings.Split(lonDeg,",")
	fmt.Print(lonDeg)
	val1,_ = strconv.ParseFloat(strings.TrimSpace(bits[0]),64)
	val2,_ = strconv.ParseFloat(strings.TrimSpace(bits[1]),64)
	val3,_ = strconv.ParseFloat(strings.TrimSpace(bits[2]),64)
	lonDec:= val1 + (val2/60) + (val3/3600)
	if "S" == lonFlag || "W" == lonFlag{
		lonDec = lonDec * -1
	}

	retFloat[0] = lonDec
	retFloat[1] = latDec

	return retFloat

}


func executor(jobs [] Worker, con int ){
	c := make(chan int);
	l:= len(jobs)
	if l < con{
		for _,w := range jobs{
			go w(c);
		}
	}else{
		s := jobs[:con]
		for _,w := range s{
			go w(c);
		}

		done :=con;


		for{
			done-= <-c
			if done == 0{
				break;
			}
		}

		executor(jobs[con:],con)
	}



}



func createThumb(filepath string, filename string) (string,error){
	// open "test.jpg"
	InfoLog.Println("opening " + filepath)
	file, err := os.Open(filepath)
	LogOnError(err, "failed to open img " + filepath)

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)


	if err !=nil{
		LogOnError(err, "failed decode")
		return "",err
	}
	defer file.Close()
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(350, 0, img, resize.Bicubic)
	thumbPath := CONF.GetThumbNailDir() + "/" + filename
	out, err := os.Create(thumbPath)

    LogOnError(err, "failed to write out thumbnail " + thumbPath)

	InfoLog.Println(" created img  " + thumbPath)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return thumbPath,err
}

func parseDate(dateString string)time.Time{
	time,err := time.Parse("2006:01:02 15:04:05",dateString)
	if err != nil{
		LogOnError(err,"failed to parse time")
	}
	return time;
}
