package api

import (
	"net/http"
	"github.com/maleck13/photoProcessor/model"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"github.com/maleck13/photoProcessor/conf"
	"strconv"
	"time"
)


func GetPicture(wr http.ResponseWriter, req *http.Request){

	vars := mux.Vars(req)
	fmt.Println("GetPicture %s", req.URL.Query())
	file := vars["file"]
	user := vars["user"]
	pic:=&model.Picture{}

	err,_ :=pic.FindByNameAndUser(file,user)

	wr.Header().Set("Content-Type", "image/jpeg")
	if nil != err{
		wr.WriteHeader(500)
		json.NewEncoder(wr).Encode(err)
	}else {
		//get file path stream back to client io.Copy(w, resp.Body)
		path := conf.CONF.GetPhotoDir() + "/" + user + "/thumbs/" + file
		f,err := os.Open(path)
		if nil != err{
			fmt.Println("err opening file " + err.Error())
		}
		io.Copy(wr,f)
	}
}

func GetYearRange(wr http.ResponseWriter, req *http.Request){

	fmt.Println("vars " ,req.URL.Query().Get("user"))
	user := req.URL.Query().Get("user")
	pic:= &model.Picture{};
	err,years := pic.GetPictureDateRange(user);
	wr.Header().Set("Content-Type", "application/json")
	if nil != err{
		wr.WriteHeader(500)
		json.NewEncoder(wr).Encode(err)
	}else {
		//get file path stream back to client io.Copy(w, resp.Body)
		json.NewEncoder(wr).Encode(years)
	}
}

func GetPicturesInRange(wr http.ResponseWriter, req * http.Request){
	user := req.URL.Query().Get("user")
	from := req.URL.Query().Get("from")
	to := req.URL.Query().Get("to")
	wr.Header().Set("Content-Type", "application/json")

	if "" == user{
		wr.WriteHeader(400)
		err := "{\"error\":\"no user specified\"}"
		json.NewEncoder(wr).Encode(err)
	}else if from == ""{
		wr.WriteHeader(400)
		err := "{\"error\":\"no from time specified\"}"
		json.NewEncoder(wr).Encode(err)
	}else if to == ""{

		wr.WriteHeader(400)
		err := "{\"error\":\"no to time specified\"}"
		json.NewEncoder(wr).Encode(err)

	}else {

		fromInt, err := strconv.ParseInt(from, 10, 64)
		toInt, err := strconv.ParseInt(to, 10, 64)

		if err != nil {
			fmt.Println("error parsing int " + err.Error())
		}

		fromTime := time.Unix(fromInt,0)
		toTime := time.Unix(toInt,0)
		fromYear := fromTime.Year()
		toYear := toTime.Year()
		fromMonth := fromTime.Month()
		toMonth := toTime.Month()
		numYears := toYear - fromYear;
		days := daysIn(toMonth, fromYear);
		pic := &model.Picture{};
		pics := make([]model.Picture, 0);

		fmt.Println("max days is ", days, numYears, fromYear,toYear, fromInt,toInt)

		for i := 0; i <= numYears; i++ {
			cYear := fromYear + i;
			fd := time.Date(cYear, toMonth, days, 0, 0, 0, 0, time.UTC)
			td := time.Date(cYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
			fmt.Println("in range cyear is ", cYear, fd.Unix(), td.Unix())
			err, retPics := pic.GetPicturesInRange(user, fd.Unix(), td.Unix())
			if nil == err {
				for _, gotPic := range retPics {
					pics = append(pics, gotPic)
				}
			}
		}

		json.NewEncoder(wr).Encode(pics)
	}




}

func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}