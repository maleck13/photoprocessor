package processor

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gosexy/exif"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"github.com/maleck13/photoProcessor/errorHandler"
	"github.com/maleck13/photoProcessor/model"
	"github.com/maleck13/photoProcessor/conf"
	"github.com/maleck13/photoProcessor/logger"
	"github.com/maleck13/photoProcessor/storage"
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
	Save() error
}

func ProcessPhotoDir(dir, user string) {
	fInfo, err := ioutil.ReadDir(dir)
	errorHandler.FailOnError(err, " Error reading dir:")

	// takes filinfo and return a func to be executed uses chan to increment a count
	buildWorker := func(f os.FileInfo) func(chan int) {
		return func(c chan int) {

			if !f.IsDir() && strings.Contains(f.Name(), ".JPG") {
				fmt.Println("processing " + f.Name())
				uc := make(chan string)
				go logMessages(uc)
				ProcessImg(f.Name(), model.Picture{}, user, uc, "internal")

			} else {
				fmt.Println("	unable to proces " + f.Name())
			}
			c <- 1
		}
	}

	//build a slice of Worker jobs to be executed
	jobs := make([]Worker, len(fInfo))

	for idx, f := range fInfo {
		jobs[idx] = buildWorker(f)
	}

	logger.InfoLog.Printf("qued up %d \n", len(jobs))

	//execute the jobs ten at a time in a go routine
	if len(jobs) > 0 {
		go executor(jobs, conf.CONF.GetConcurrentJobs())
	} else {
		logger.InfoLog.Println("no files found in dir")
	}

}

func logMessages(messages chan string) {
	for m := range messages {
		fmt.Println(" message update " + m)
	}
}

func CreateMessage(message, status, jobid string) string {
	msg := model.UPDATE_MESSAGE{message, status, jobid, "PICTURE"}
	json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error " + err.Error())
	}
	return string(json)
}

func ProcessImg(fileName string, pic model.Picture, user string, updateChanel chan string, jobId string) {

	msg := CreateMessage("starting processing img ", "pending", jobId)
	fmt.Println("made message " + msg)
	updateChanel <- msg
	defer close(updateChanel)

	reader := exif.New()

	path := conf.CONF.GetPhotoDir() + "/" + user + "/" + fileName
	completedPath := conf.CONF.GetPhotoDir() + "/" + user + "/completed/" + fileName

	fmt.Println("Path to img " + path)
	fmt.Println("completed path " + completedPath)

	err := reader.Open(path)
	errorHandler.LogOnError(err, "Error reading data from "+path)
	if nil != err {
		msg = CreateMessage("Error reading data from "+path+" "+err.Error(), "error", jobId)
		updateChanel <- msg
		return
	}

	tags := reader.Tags


	for key, value := range tags {
		fmt.Println("Key:", key, "Value:", value)
	}


	var lonLat []float64

	err = validateLonLat(tags)



	var hasLonLat, hasTime bool;
	if err != nil {
		msg = CreateMessage("Error no longlat exif data ", "error", jobId)
		errorHandler.LogOnError(err, "missing lonlat data")
		hasLonLat = false
	} else {
		lonLat = convertDegToDec(tags[LATITUDE], tags[NORTH_OR_SOUTH_LAT], tags[LONGITUDE], tags[EAST_OR_WEST_LON])
		hasLonLat = true;
	}

	err = validateTime(tags)
	if err != nil {
		msg = CreateMessage("Error no time exif data ", "error", jobId)
		updateChanel <- msg
		errorHandler.LogOnError(err, "missing data")
		hasTime = false;
	}else{
		hasTime = true;
	}


	pic.Complete = (hasTime && hasLonLat)
	pic.LonLat = lonLat
	pic.Name = fileName
	pic.Img = fileName
	pic.Path = completedPath
	pic.Tags = []string{};
	thumb, err := createThumb(path, fileName, user, tags)
	if err != nil {
		errorHandler.LogOnError(err, "failed to create thumb ignoring img "+fileName)
	}
	msg = CreateMessage("Thumbnail created  "+thumb, "pending", jobId)
	updateChanel <- msg

	date := parseDate(tags[DATE_TIME_KEY])

	pic.Thumb = thumb
	pic.Time = date
	pic.Year = date.String()[0:4]
	pic.User = user
	pic.TimeStamp = date.Unix()
	logger.InfoLog.Println(pic)
	err = pic.Save()
	msg = CreateMessage("Saved to db ", "complete", jobId)
	updateChanel <- msg
	if err != nil {
		msg = CreateMessage("failed Save to db ", "pending", jobId)
		errorHandler.LogOnError(err, "failed to save picture")
		//move to failed dir
	}


}

func ReadExifData(filePath string) (map[string]string, error) {
	reader := exif.New()
	_, err := os.Stat(filePath)
	if nil != err {
		return nil, err
	}

	err = reader.Open(filePath)
	//LogOnError(err, "failed to open "+filePath)
	if nil != err {
		return nil, err
	}

	tags := reader.Tags

	return tags, nil
}

func validateLonLat(info map[string]string) error {
	ok, hasKey := info[LONGITUDE]
	fmt.Println("validate lon lat val is " + ok)
	if !hasKey || "" == ok{
		return errors.New("no " + LONGITUDE + " field")
	}
	ok, hasKey = info[LATITUDE]

	if !hasKey || "" == ok{
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
		latDec = latDec * -1
	}

	bits = strings.Split(lonDeg, ",")
	fmt.Print(lonDeg)
	val1, _ = strconv.ParseFloat(strings.TrimSpace(bits[0]), 64)
	val2, _ = strconv.ParseFloat(strings.TrimSpace(bits[1]), 64)
	val3, _ = strconv.ParseFloat(strings.TrimSpace(bits[2]), 64)
	lonDec := val1 + (val2 / 60) + (val3 / 3600)
	if "S" == lonFlag || "W" == lonFlag {
		lonDec = lonDec * -1
	}

	retFloat[0] = lonDec
	retFloat[1] = latDec

	return retFloat

}

func executor(jobs []Worker, con int) {
	c := make(chan int)
	l := len(jobs)
	if l < con {
		for _, w := range jobs {
			go w(c)
		}
	} else {
		s := jobs[:con]
		for _, w := range s {
			go w(c)
		}

		done := con

		for {
			done -= <-c
			if done == 0 {
				break
			}
		}

		executor(jobs[con:], con)
	}

}

func createThumb(filepath string, filename string, user string, exif map[string]string) (string, error) {
	// open "test.jpg"
	logger.InfoLog.Println("opening " + filepath)
	file, err := os.Open(filepath)
	if err != nil {
		errorHandler.LogOnError(err, "failed to open filepath "+filepath)
	}
	defer file.Close()
	thumbPath := conf.CONF.GetPhotoDir() + "/" + user + "/thumbs/" + filename

	fmt.Println("thumbpath is " + thumbPath)

	errorHandler.LogOnError(err, "failed to open img "+filepath)

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)

	if err != nil {
		// just move the img to thumbs
		file, err := os.Open(filepath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		errorHandler.LogOnError(err, "failed decode")
		logger.InfoLog.Println(" error decoding copying file")
		fc, err := os.Create(thumbPath)
		if err != nil {
			return "", err
		}
		logger.InfoLog.Println(" created new path " + thumbPath)
		defer fc.Close()

		_, err = io.Copy(fc, file)

		if err != nil {
			// just move the img to thumbs
			errorHandler.LogOnError(err, "failed to copy file")
		}

		if conf.CONF.GetAwsEnabled() {
			logger.InfoLog.Println(" AWS ENABLED *** addding thub to aws *** ")
			thumbPath, err = storage.PutInBucket(thumbPath, "/"+user+"/thumbs/"+filename)
			if nil != err {
				logger.ErrorLog.Println(" failed to add to aws " + err.Error())
			}
		}

		return thumbPath, err
	}




	orientate := exif["Orientation"]
	if "Right-top" == orientate {
		logger.InfoLog.Println("rotating img 270 degrees")
		img = imaging.Rotate270(img)
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
//	var percentHeight, percentWidth int
//	percentHeight = (img.Bounds().Max.Y / 100) * 30
//	percentWidth = (img.Bounds().Max.X / 100) * 15

	m := imaging.Thumbnail(img, 300, 300,imaging.Lanczos)

	out, err := os.Create(thumbPath)

	errorHandler.LogOnError(err, "failed to write out thumbnail "+thumbPath)

	logger.InfoLog.Println(" created img  " + thumbPath)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	if conf.CONF.GetAwsEnabled() {
		logger.InfoLog.Println(" AWS ENABLED *** addding thub to aws *** ")
		thumbPath, err = storage.PutInBucket(thumbPath, "/"+user+"/thumbs/"+filename)
		if nil != err {
			logger.ErrorLog.Println(" failed to add to aws " + err.Error())
		}
	}
	return thumbPath, err
}

func parseDate(dateString string) time.Time {
	time, err := time.Parse("2006:01:02 15:04:05", dateString)
	if err != nil {
		errorHandler.LogOnError(err, "failed to parse time")
	}
	return time
}
