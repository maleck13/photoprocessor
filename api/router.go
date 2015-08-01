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

	router.HandleFunc("/picture/{id}",UpdatePictureData).Methods("PUT")
	router.HandleFunc("/health",Ping).Methods("GET")

	sub := router.PathPrefix("/pictures").Subrouter()
	sub.Path("/").HandlerFunc(GetPicturesInRange).Methods("GET")
	sub.HandleFunc("/{user}/{file}",GetPicture).Methods("GET")
	sub.HandleFunc("/{user}/{file}",DeletePicture).Methods("DELETE")
	sub.HandleFunc("/{user}/upload",UploadHandler).Methods("POST")
	sub.Path("/range").HandlerFunc(GetYearRange).Methods("GET")
	sub.Path("/incomplete").HandlerFunc(GetIncompletePictures).Methods("GET")
	sub.Path("/byloc").HandlerFunc(GetPicturesByLocation).Methods("GET")
	sub.Path("/bytag").HandlerFunc(GetPicturesByTag).Methods("GET")


	http.Handle("/", router)

	return router;
}


