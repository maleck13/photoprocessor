package api

import (
	"net/http"
	"github.com/maleck13/photoProcessor/model"
	"encoding/json"
)



func IndexRoute(wr http.ResponseWriter, req *http.Request){
pic := model.Picture{}
json.NewEncoder(wr).Encode(pic)
}
