package handler

import (
    "net/http"
    "encoding/json"

    "datasource"
)

func GetLogs(w http.ResponseWriter, r *http.Request) {
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

    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("timestamp", ">", tsRange.Begin))
    conditions = append(conditions, datasource.NewCondition("timestamp", "<", tsRange.End))

    rl := datasource.ListRequestLog(db, conditions, nil, nil)
    json.NewEncoder(w).Encode(rl)
}

func GetServiceLogs(w http.ResponseWriter, r *http.Request) {
    // TBD
}

func GetRequestTypeLogs(w http.ResponseWriter, r *http.Request) {
    // TBD
}

func GetRequestLogs(w http.ResponseWriter, r *http.Request) {
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

    uuid := getUuid(r)

    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("uuid", "=", uuid))
    conditions = append(conditions, datasource.NewCondition("timestamp", ">", tsRange.Begin))
    conditions = append(conditions, datasource.NewCondition("timestamp", "<", tsRange.End))

    rl := datasource.ListRequestLog(db, conditions, nil, nil)
    json.NewEncoder(w).Encode(rl)
}
