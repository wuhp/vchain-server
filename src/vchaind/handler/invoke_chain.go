package handler

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "vchaind/model"
)

func GetAllInvokeChains(w http.ResponseWriter, r *http.Request) {
    ivkchains := model.ListInvokeChain(nil, nil, nil)
    json.NewEncoder(w).Encode(ivkchains)    
}

func GetInvokeChains(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    service := vars["service"]
    category := vars["category"]

    st := model.RequestType{
        Service: service,
        Category: category,
    }

    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("header", "=", model.RequestType2string(&st)))
    ivkchains := model.ListInvokeChain(conditions, nil, nil)
    json.NewEncoder(w).Encode(ivkchains)
}

func GetInvokeChain(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    service := vars["service"]
    category := vars["category"]
    id := vars["id"]

    st := model.RequestType{
        Service: service,
        Category: category,
    }

    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("header", "=", model.RequestType2string(&st)))
    conditions = append(conditions, model.NewCondition("id", "=", id))

    ivkchains := model.ListInvokeChain(conditions, nil, nil)
    if len(ivkchains) == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(ivkchains[0])
}

func GetInvokeChainHeaderRequests(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    service := vars["service"]
    category := vars["category"]
    id := vars["id"]
    
    st := model.RequestType{
        Service: service,
        Category: category,
    }   
    
    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("header", "=", model.RequestType2string(&st)))
    conditions = append(conditions, model.NewCondition("id", "=", id))
    
    ivkchains := model.ListInvokeChain(conditions, nil, nil)
    if len(ivkchains) == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    uuids := make([]string, 0)
    conditions = make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("invoke_chain_id", "=", id))

    rgs := model.ListRequestGroup(conditions, nil, nil)
    for _, rg := range rgs {
        uuids = append(uuids, rg.Uuid)
    }

    json.NewEncoder(w).Encode(uuids)
}
