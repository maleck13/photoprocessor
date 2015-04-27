package api

import (
	"net/http"
	"log"
)



func StartApi(){

	NewApiRouter()
	log.Fatal(http.ListenAndServe(":9002", nil))
}
