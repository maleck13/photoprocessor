package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewApiRouter() *mux.Router{

	router := mux.NewRouter()
	router.HandleFunc("/pictures/{user}/{file}",GetPicture).Methods("GET")
	router.HandleFunc("/pictures/{user}/upload",UploadHandler).Methods("POST")
	router.HandleFunc("/pictures/range",GetYearRange).Methods("GET")
	router.HandleFunc("/pictures",GetPicturesInRange).Methods("GET")
	router.HandleFunc("/health",Ping).Methods("GET")
	http.Handle("/", router)

	return router;
}


