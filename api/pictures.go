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
		path := conf.CONF.GetPhotoDir() + "/" + user + "/" + file
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
