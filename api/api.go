package api

import (
	"net/http"
	"log"
)

type ApiError struct{
	Error string
}


func StartApi(){

	NewApiRouter()
	log.Fatal(http.ListenAndServe(":9002", nil))
}
