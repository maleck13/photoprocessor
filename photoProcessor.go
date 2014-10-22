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

	EAST_OR_WEST_LON   = "East or West Longitude"
	NORTH_OR_SOUTH_LAT = "North or South Latitude"
	DATE_TIME_KEY      = "Date and Time (Original)"
	LONGITUDE          = "Longitude"
	LATITUDE           = "Latitude"
)

type Worker func(chan int)

type Persister interface {
	Save()error
}

func ProcessPhotoDir() {

	fInfo, err := ioutil.ReadDir(CONF.GetPhotoDir())
	FailOnError(err, " Error reading dir:")

	// takes filinfo and return a func to be executed uses chan to increment a count
	buildWorker := func(f os.FileInfo) func(chan int) {
		return func(c chan int) {

			if !f.IsDir() && strings.Contains(f.Name(), ".JPG") {
				fmt.Println("processing " + f.Name())
				ProcessImg(f.Name(),Picture{},CONF)
			}
			c <- 1
		}
	}

	//build a slice of Worker jobs to be executed
	jobs := make([]Worker, len(fInfo))

	for idx, f := range fInfo {
		jobs[idx] = buildWorker(f)
	}

	InfoLog.Printf("qued up %d \n", len(jobs))

	//execute the jobs ten at a time in a go routine
	if len(jobs) > 0 {
	 go executor(jobs, CONF.GetConcurrentJobs())
	}else{
		InfoLog.Println("no files found in dir")
	}

}

func ProcessImg(fileName string, pic Picture, conf *CONFIG) {
	reader := exif.New()
	path := conf.GetPhotoDir() + "/" + fileName
	completedPath := conf.GetProcessedPhotoDir() + "/" + fileName
	err := reader.Open(path)

	LogOnError(err, "failed to open "+path)

	tags := reader.Tags
	var lonLat []float64

	err = validateLonLat(tags);

	if err != nil {
		if conf.GetUseDefaultLonLat(){
			InfoLog.Println("using default lon lat ")
			lonLat = conf.GetDefaultLonLat()
		}else {
			LogOnError(err, "missing data")
			return
		}
	}else {
		lonLat = convertDegToDec(tags[LATITUDE], tags[NORTH_OR_SOUTH_LAT], tags[LONGITUDE], tags[EAST_OR_WEST_LON])
	}

	err = validateTime(tags);
	if err != nil {
		LogOnError(err, "missing data")
		return
	}


	pic.LonLat = lonLat
	pic.Name = fileName
	pic.Path = completedPath
	thumb, err := createThumb(path, fileName, conf)
	date := parseDate(tags[DATE_TIME_KEY])

	if err != nil {
		LogOnError(err, "failed to create thumb ignoring img "+fileName)
		return
	}
	pic.Thumb = thumb
	pic.Time = date
	pic.TimeStamp = date.Unix()
	err = pic.Save()

	if err != nil {
		LogOnError(err, "failed to save picture")
	}

	err = copyAndRemove(fileName, conf)
	if err != nil {
		FailOnError(err, "failed to save picture")
	}


}

func copyAndRemove (fileName string, conf * CONFIG) error{
	path := conf.GetPhotoDir() + "/" + fileName
	completedPath := CONF.GetProcessedPhotoDir() + "/" +  fileName

	dir,err := os.Stat(CONF.GetPhotoDir())
	if err != nil {
		return err
	}

	if dir.IsDir() == false{
		return errors.New("photo dir is not a dir ")
	}

	dir,err = os.Stat(conf.GetProcessedPhotoDir())

	if err != nil {
		return err
	}

	if dir.IsDir() == false{
		return errors.New(" completed photo dir is not a dir ")
	}


	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()


	fc, err := os.Create(completedPath)
	if err != nil {
		return err
	}
	defer fc.Close()

	_, err = io.Copy(fc, f)

	if err != nil {
		return err
	}

	err = os.Remove(path)

	return err
}

func validateLonLat(info map[string]string) error {
	_, ok := info[LONGITUDE]
	if !ok {
		return errors.New("no " + LONGITUDE + " field")
	}
	_, ok = info[LATITUDE]

	if !ok {
		return errors.New("no " + LATITUDE + " field")
	}
	return nil
}

func validateTime(info map[string]string) error {
	_, ok := info[DATE_TIME_KEY]
	if !ok {
		return errors.New("no " + DATE_TIME_KEY + " present")
	}
	return nil
}

func convertDegToDec(latDeg string, latFlag string, lonDeg string, lonFlag string) []float64 {
	//12° 34" 56' = 12 + (34/60) + (56/3600) = 12.582222222222222222222
	bits := strings.Split(latDeg, ",")

	fmt.Print(bits[0] + " : " + bits[1] + " : " + bits[2])
	val1, _ := strconv.ParseFloat(strings.TrimSpace(bits[0]), 64)
	val2, _ := strconv.ParseFloat(strings.TrimSpace(bits[1]), 64)
	val3, _ := strconv.ParseFloat(strings.TrimSpace(bits[2]), 64)
	retFloat := make([]float64, 2)
	fmt.Printf(" float vals %f %f %f ", val1, val2, val3)
	latDec := val1 + (val2 / 60) + (val3 / 3600)

	if "N" != latFlag && "E" != latFlag {
		latDec = latDec*-1
	}

	bits = strings.Split(lonDeg, ",")
	fmt.Print(lonDeg)
	val1, _ = strconv.ParseFloat(strings.TrimSpace(bits[0]), 64)
	val2, _ = strconv.ParseFloat(strings.TrimSpace(bits[1]), 64)
	val3, _ = strconv.ParseFloat(strings.TrimSpace(bits[2]), 64)
	lonDec := val1 + (val2 / 60) + (val3 / 3600)
	if "S" == lonFlag || "W" == lonFlag {
		lonDec = lonDec*-1
	}

	retFloat[0] = lonDec
	retFloat[1] = latDec

	return retFloat

}


func executor(jobs [] Worker, con int) {
	c := make(chan int);
	l := len(jobs)
	if l < con {
		for _, w := range jobs {
			go w(c);
		}
	}else {
		s := jobs[:con]
		for _, w := range s {
			go w(c);
		}

		done := con;


		for {
			done-= <-c
			if done == 0 {
				break;
			}
		}

		executor(jobs[con:], con)
	}

}



func createThumb(filepath string, filename string, conf * CONFIG) (string, error) {
	// open "test.jpg"
	InfoLog.Println("opening " + filepath)
	file, err := os.Open(filepath)
	LogOnError(err, "failed to open img "+filepath)

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)


	if err != nil {
		LogOnError(err, "failed decode")
		return "", err
	}
	defer file.Close()
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(350, 0, img, resize.Bicubic)
	thumbPath := conf.GetThumbNailDir() + "/" + filename
	out, err := os.Create(thumbPath)

	LogOnError(err, "failed to write out thumbnail "+thumbPath)

	InfoLog.Println(" created img  " + thumbPath)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return thumbPath, err
}

func parseDate(dateString string) time.Time {
	time, err := time.Parse("2006:01:02 15:04:05", dateString)
	if err != nil {
		LogOnError(err, "failed to parse time")
	}
	return time;
}
