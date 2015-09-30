package handler

import (
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

    reqs := make([]*datasource.Request, 0)
    if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, req := range reqs {
        req.Save(db)
    }
}

func PostRequestLog(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }    
    
    rlogs := make([]*datasource.RequestLog, 0)
    if err := json.NewDecoder(r.Body).Decode(&rlogs); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, rlog := range rlogs {
        rlog.Save(db)
    }
}
