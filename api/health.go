package api
import (
"net/http"
"encoding/json"
)

func Ping(wr http.ResponseWriter, req *http.Request){
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode("{\"health\":\"ok\"}")
}
