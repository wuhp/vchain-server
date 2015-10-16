package handler

import (
    "time"
    "log"
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
        if datasource.GetRequest(db, req.Uuid) != nil {
            log.Printf("INFO: request %s exists, skip to save ...\n", req.Uuid)
            continue
        }
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
        if datasource.GetRequestLog(db, rlog.Uuid, rlog.Timestamp) != nil {
            log.Printf("INFO: request log `%s %d` exists, skip to save ...\n", rlog.Uuid, rlog.Timestamp)
            continue
        }
        rlog.Save(db)
    }
}
