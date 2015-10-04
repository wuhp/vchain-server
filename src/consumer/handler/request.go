package handler

import (
    "net/http"
    "encoding/json"

    "datasource"
)

func GetRequests(w http.ResponseWriter, r *http.Request) {
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
    conditions = append(conditions, datasource.NewCondition("begin_ts", ">", tsRange.Begin))
    conditions = append(conditions, datasource.NewCondition("begin_ts", "<", tsRange.End))

    rs := datasource.ListRequest(db, conditions, nil, nil)
    json.NewEncoder(w).Encode(rs)
}

func GetRequest(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    if req == nil {
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(req)
}

func GetRequestInvokeChain(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    switch {
    case req == nil:
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    case req.GroupUuid == "":
        http.Error(w, "Request has not processed yet", http.StatusBadRequest)
        return
    }

    rg := datasource.GetRequestGroup(db, req.GroupUuid)
    ic := datasource.GetInvokeChain(db, rg.InvokeChainId)
    json.NewEncoder(w).Encode(ic)
}

func GetRequestRequestGroup(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    switch {
    case req == nil:
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    case req.GroupUuid == "":
        http.Error(w, "Request has not processed yet", http.StatusBadRequest)
        return
    }

    rg := datasource.GetRequestGroup(db, req.GroupUuid)

    result := struct {
        Requests     []*datasource.Request `json:"requests"`
        ParentsIndex []int                 `json:"parents_index"`
    }{}

    result.ParentsIndex = rg.ParentsIndex
    result.Requests = rg.DetailRequests(db)

    json.NewEncoder(w).Encode(result)
}

func GetRequestRootRequest(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    switch {
    case req == nil:
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    case req.GroupUuid == "":
        http.Error(w, "Request has not processed yet", http.StatusBadRequest)
        return
    }

    root := datasource.GetRequest(db, req.GroupUuid)
    json.NewEncoder(w).Encode(root)
}

func GetRequestParent(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    if req == nil {
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    }

    parent := datasource.GetRequest(db, req.ParentUuid)
    if parent == nil {
        http.Error(w, "Request has no parent", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(parent)
}

func GetRequestChildren(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    req := datasource.GetRequest(db, getUuid(r))
    if req == nil {
        http.Error(w, "Invalid request uuid", http.StatusNotFound)
        return
    }

    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("parent_uuid", "=", req.Uuid))
    rs := datasource.ListRequest(db, conditions, nil, nil)
    json.NewEncoder(w).Encode(rs)
}
