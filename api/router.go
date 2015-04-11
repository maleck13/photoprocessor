package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"github.com/maleck13/photoProcessor/model"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewApiRouter() *mux.Router{

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/",IndexRoute);


	return router;
}

func IndexRoute(wr http.ResponseWriter, req *http.Request){
	pic := model.Picture{}
	json.NewEncoder(wr).Encode(pic)
}
