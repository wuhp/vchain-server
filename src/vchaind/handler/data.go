package handler

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "vchaind/model"
)

var (
    ServiceStartEvent = "service.start"
    ServiceStopEvent  = "service.stop"
    RequestBeginEvent = "request.begin"
    RequestEndEvent   = "request.end"
)

type Data struct {
    Event   string                 `json:"event"`
    Payload map[string]interface{} `json:"payload"`
}

func getAppId(secret string) (int, bool) {
    // TBD
    // Internal error -> panic
    // secret valid return app_id, true
    // secret invalid return 0, false
    return 1, true
}

func PostData(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    secret := vars["secret"]

    _, exist := getAppId(secret)
    if !exist {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    in := make([]Data, 0)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, data := range in {
        switch(data.Event) {
        case ServiceStartEvent:
            s := new(model.Service)
            s.Uuid = data.Payload["uuid"].(string)
            s.Category = data.Payload["category"].(string)
            s.Instance = data.Payload["datastance"].(string)
            s.Hostname = data.Payload["hostname"].(string)
            s.StartTs = data.Payload["start_ts"].(int64)
            s.Save()
        case ServiceStopEvent:
            conditions := make([]*model.Condition, 0)
            conditions = append(conditions, model.NewCondition("uuid", "=", data.Payload["uuid"]))
            s := model.ListService(conditions, nil, nil)[0]
            s.StopTs = data.Payload["stop_ts"].(int64)
            s.Update()
        case RequestBeginEvent:
            r := new(model.Request)
            r.Uuid = data.Payload["uuid"].(string)
            r.ServiceUuid = data.Payload["service_uuid"].(string)
            r.ParentUuid = data.Payload["parent_uuid"].(string)
            r.Category = data.Payload["category"].(string)
            r.BeginTs = data.Payload["begin_ts"].(int64)
            bs, _ := json.Marshal(data.Payload["begin_metadata"])
            r.BeginMetadata = string(bs)
            r.Save()
        case RequestEndEvent:
            conditions := make([]*model.Condition, 0)
            conditions = append(conditions, model.NewCondition("uuid", "=", data.Payload["uuid"]))
            r := model.ListRequest(conditions, nil, nil)[0]
            r.EndTs = data.Payload["end_ts"].(int64)
            bs, _ := json.Marshal(data.Payload["end_metadata"])
            r.EndMetadata = string(bs)
            r.Update()
        default:
            http.Error(w, fmt.Sprintf("Invalid event %s", data.Event), http.StatusBadRequest)
            return
        }
    }
}
