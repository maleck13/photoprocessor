package api

import (
	"net/http"
	"log"
)



func StartApi(){

	router:= NewApiRouter()
	log.Fatal(http.ListenAndServe(":8881", router))
}
