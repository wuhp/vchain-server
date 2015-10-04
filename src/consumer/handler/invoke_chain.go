package handler

import (
    "net/http"
    "encoding/json"

    "datasource"
)

func GetAllInvokeChains(w http.ResponseWriter, r *http.Request) {
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

    ivkchains := datasource.QueryInvokeChain(db, tsRange, nil)
    json.NewEncoder(w).Encode(ivkchains)
}

func GetInvokeChains(w http.ResponseWriter, r *http.Request) {
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

    rt := getRequestType(r)

    ivkchains := datasource.QueryInvokeChain(db, tsRange, rt)
    json.NewEncoder(w).Encode(ivkchains)
}

func GetInvokeChain(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    rt := getRequestType(r)
    id := getId(r)

    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("header", "=", datasource.RequestType2string(rt)))
    conditions = append(conditions, datasource.NewCondition("id", "=", id))

    ivkchains := datasource.ListInvokeChain(db, conditions, nil, nil)
    if len(ivkchains) == 0 {
        http.Error(w, "Invalid invoke chain (service, category, id)", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(ivkchains[0])
}

func GetInvokeChainRootRequests(w http.ResponseWriter, r *http.Request) {
    db := getDb(r)
    if db == nil {
        http.Error(w, "Invalid project id", http.StatusNotFound)
        return
    }
    defer db.Close()

    rt := getRequestType(r)
    id := getId(r)

    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("header", "=", datasource.RequestType2string(rt)))
    conditions = append(conditions, datasource.NewCondition("id", "=", id))

    ivkchains := datasource.ListInvokeChain(db, conditions, nil, nil)
    if len(ivkchains) == 0 {
        http.Error(w, "Invalid invoke chain (service, category, id)", http.StatusNotFound)
        return
    }

    rs := datasource.FindRequestsByInvokeChain(db, ivkchains[0].Id)
    json.NewEncoder(w).Encode(rs)
}
