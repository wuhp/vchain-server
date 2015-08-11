package handler

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "vchaind/model"
)

func GetRequestOverview(w http.ResponseWriter, r *http.Request) {
    // TBD
}

func GetRequestTypes(w http.ResponseWriter, r *http.Request) {
    rts := model.GetRequestTypes()
    json.NewEncoder(w).Encode(rts)
}

func GetRequests(w http.ResponseWriter, r *http.Request) {
    rs := model.ListRequest(nil, nil, nil)
    json.NewEncoder(w).Encode(rs)
}

func GetRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    req := model.GetRequest(uuid)
    if req == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(req)
}

func GetRequestInvokeChain(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    req := model.GetRequest(uuid)
    if req == nil || req.GroupUuid == "" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    rg := model.GetRequestGroup(req.GroupUuid)
    if req == nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    ic := model.GetInvokeChain(rg.InvokeChainId)
    if ic == nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(ic)
}

func GetRequestRequestGroup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    req := model.GetRequest(uuid)
    if req == nil || req.GroupUuid == "" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    rg := model.GetRequestGroup(req.GroupUuid)
    if req == nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    result := struct {
        Requests     []*model.Request `json:"requests"`
        ParentsIndex []int            `json:"parents_index"`
    }{}

    result.ParentsIndex = rg.ParentsIndex
    result.Requests = rg.DetailRequests()

    json.NewEncoder(w).Encode(result)
}

func GetRequestRootRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    req := model.GetRequest(uuid)
    if req == nil || req.GroupUuid == "" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    root := model.GetRequest(req.GroupUuid)
    if root == nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(root)
} 

func GetRequestParent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    req := model.GetRequest(uuid)
    if req == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }    

    parent := model.GetRequest(req.ParentUuid)
    json.NewEncoder(w).Encode(parent)
}

func GetRequestChildren(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]
    
    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("parent_uuid", "=", uuid))
    rs := model.ListRequest(conditions, nil, nil)
    json.NewEncoder(w).Encode(rs)
}
