package handler

import (
    "net/http"
    "encoding/json"

    "datasource"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }

    tsRange := getTimeRange(r)
    if tsRange == nil {
        http.Error(w, "Invalid time range value", http.StatusBadRequest)
        return
    }

    services := datasource.GetServices(db, tsRange)
    json.NewEncoder(w).Encode(services)
}

func GetServiceChain(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }

    tsRange := getTimeRange(r)
    if tsRange == nil {
        http.Error(w, "Invalid time range value", http.StatusBadRequest)
        return
    }

    pairs := datasource.GetServiceChain(db, tsRange)
    json.NewEncoder(w).Encode(pairs)
}
