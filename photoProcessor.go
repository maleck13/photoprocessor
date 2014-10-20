package main

import (
	"fmt"
	"github.com/gosexy/exif"
	"io/ioutil"
	"strings"
	"os"
	"strconv"
	"errors"
)

const (
	PIC_DIR = "/Users/craigbrookes/Pictures"
	EAST_OR_WEST_LON = "East or West Longitude"
	NORTH_OR_SOUTH_LAT = "North or South Latitude"
	LONGITUDE = "Longitude"
	LATITUDE = "Latitude"
)


type Worker func(chan int)

func ProcessPhotoDir() {
	fInfo, err := ioutil.ReadDir(PIC_DIR)
	if err != nil {
		fmt.Printf("Error reading dir: %s", err.Error())
	}


   execute := func (f os.FileInfo ) func(chan int){
	 return func (c chan int){

	   if !f.IsDir() && strings.Contains(f.Name(), ".JPG") {
		   ProcessImg(f.Name())

	   }
	   c <- 1

	 }
   }


	jobs:=make([]Worker, len(fInfo))
	for idx, f := range fInfo {
		jobs[idx] = execute(f)
	}

	go executor(jobs,10)


}

func ProcessImg(filePath string ) {
	//fmt.Println("process img", filePath);
	reader := exif.New()
	path := PIC_DIR + "/" + filePath
	err := reader.Open(path)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	pic := Picture{}
	tags := reader.Tags
	err = validateLonLat(tags)
	if err != nil{
				fmt.Printf("error missing lonLat data %s", err.Error())
	}else{
	lonLat := convertDegToDec(tags[LATITUDE],tags[NORTH_OR_SOUTH_LAT],tags[LONGITUDE],tags[EAST_OR_WEST_LON])
		pic.LonLat = lonLat
		pic.Name = filePath
		pic.Path = path
		err:=SavePic(pic)
		if err != nil{
			return
		}

		_,fPic := FindByName(filePath);

		fmt.Println(fPic.LonLat)

	}



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
