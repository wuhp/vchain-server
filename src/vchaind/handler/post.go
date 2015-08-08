package handler

import (
    "time"
    "net/http"
    "encoding/json"

    "vchaind/model"
)

var (
    RequestBeginEvent = "request.begin"
    RequestEndEvent   = "request.end"
)

type Data struct {
    Event   string        `json:"event"`
    Request model.Request `json:"payload"`
}

func PostData(w http.ResponseWriter, r *http.Request) {
    data := make([]Data, 0)
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, item := range data {
        if item.Request.Uuid == "" {
            // TBD
            continue
        }

        switch item.Event {
        case RequestBeginEvent:
            if item.Request.Exist() {
                // TBD
                continue
            }

            item.Request.GroupUuid = ""
            item.Request.CreateTs = time.Now().UTC().Unix()
            item.Request.UpdateTs = item.Request.CreateTs
            item.Request.Save()
        case RequestEndEvent:
            req := model.GetRequest(item.Request.Uuid)
            if req == nil {
                // TBD
                continue
            }

            req.EndTs = item.Request.EndTs
            req.EndMetadata = item.Request.EndMetadata
            req.UpdateTs = time.Now().UTC().Unix()
            req.Update()
        default:
            // TBD
            continue
        }
    }
}
