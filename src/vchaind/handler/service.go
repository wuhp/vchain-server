package handler

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "vchaind/model"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
    services := model.GetServices()
    json.NewEncoder(w).Encode(services)
}

func GetServiceChain(w http.ResponseWriter, r *http.Request) {
    pairs := model.GetServiceChain()
    json.NewEncoder(w).Encode(pairs)
}

func GetServiceRequestCategories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]

    categories := make([]string, 0)

    types := model.GetRequestTypes()
    for _, t := range types {
        if t.Service == name {
            categories = append(categories, t.Category)
        }
    }

    json.NewEncoder(w).Encode(categories)
}
