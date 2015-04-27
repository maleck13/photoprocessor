package api

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"github.com/maleck13/photoProcessor/logger"
	"github.com/maleck13/photoProcessor/model"
	"github.com/maleck13/photoProcessor/conf"
	"github.com/maleck13/photoProcessor/processor"
	"github.com/maleck13/photoProcessor/messaging"
	"github.com/gorilla/mux"
	"encoding/json"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	vars := mux.Vars(r)
	fmt.Println("vars ", vars , header.Filename)
	if err != nil {
		logger.ErrorLog.Println("err with upload " + err.Error())
		json.NewEncoder(w).Encode(err)
		return
	}

	user := vars["user"]


	defer file.Close()
	path := conf.CONF.GetPhotoDir() + "/" + user
	err = os.MkdirAll(path,os.ModePerm)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege " + err.Error())
		return
	}
	fullPath := path + "/" + header.Filename
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege " + err.Error())
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}




	uidKey := user + header.Filename;

	updates := make(chan string)
	messaging.SetUpResponseQue(uidKey)
	go messaging.UpdateJob( uidKey, updates)
	go processor.ProcessImg(header.Filename, model.Picture{}, user, updates, uidKey)

	json.NewEncoder(w).Encode(&model.Message{Name:header.Filename, File:fullPath, User:user, ResKey:uidKey})
}
