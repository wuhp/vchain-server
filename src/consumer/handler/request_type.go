package handler

import (
    "net/http"
    "encoding/json"

    "datasource"
)

func GetRequestTypes(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    tsRange := getTimeRange(r)
    if tsRange == nil {
        http.Error(w, "Invalid time range value", http.StatusBadRequest)
        return
    }

    rts := datasource.GetRequestTypes(db, tsRange)
    json.NewEncoder(w).Encode(rts)
}
