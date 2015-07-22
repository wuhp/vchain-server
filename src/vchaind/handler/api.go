package handler

import (
    "fmt"
    "net/http"

//    "github.com/gorilla/mux"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func PostData(w http.ResponseWriter, r *http.Request) {
    // vars := mux.Vars(r)
    // secret := vars["secret"]
    // TBD
}
