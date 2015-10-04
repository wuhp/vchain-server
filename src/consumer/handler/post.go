package handler

import (
    "time"
    "net/http"
    "encoding/json"

    "datasource"
)

func PostRequest(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    reqs := make([]*datasource.Request, 0)
    if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, req := range reqs {
        req.CreateTs = time.Now().UTC().Unix()
        req.Save(db)
    }
}

func PostRequestLog(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    rlogs := make([]*datasource.RequestLog, 0)
    if err := json.NewDecoder(r.Body).Decode(&rlogs); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, rlog := range rlogs {
        rlog.Save(db)
    }
}
