package handler

import (
    "fmt"
    "net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}
