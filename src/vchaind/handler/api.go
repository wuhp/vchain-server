package handler

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func PostData(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    secret := vars["secret"]
    // TBD
}
