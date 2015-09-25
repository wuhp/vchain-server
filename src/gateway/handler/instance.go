package handler

import (
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "gateway/model"
)

////////////////////////////////////////////////////////////////////////////////

func QueryInstance(w http.ResponseWriter, r *http.Request) {
    host := r.URL.Query().Get("host")
    port := r.URL.Query().Get("port")

    conditions := make([]*model.Condition, 0)
    if len(host) != 0 {
        conditions = append(conditions, model.NewCondition("host", "=", host))
    }

    if len(port) != 0 {
        conditions = append(conditions, model.NewCondition("port", "=", port))
    }

    json.NewEncoder(w).Encode(model.ListInstance(conditions, nil, nil))
}

func CreateInstance(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Host          string `json:"host"`
        Port          int    `json:"port"`
        AdminUser     string `json:"admin_user"`
        AdminPassword string `json:"admin_password"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "ERROR: Fail to decode json for creating instance", http.StatusBadRequest)
        return
    }

    if model.GetInstanceByHostPort(in.Host, in.Port) != nil {
        w.WriteHeader(http.StatusConflict)
        return
    }

    out := new(model.Instance)
    out.Host = in.Host
    out.Port = in.Port
    out.AdminUser = in.AdminUser
    out.AdminPassword = in.AdminPassword
    out.Active = false
    out.Save()

    json.NewEncoder(w).Encode(out)
}


func GetInstance(w http.ResponseWriter, r *http.Request) {
    instance := parseInstance(r)
    if instance == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(instance)
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
    instance := parseInstance(r)
    if instance == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    defer r.Body.Close()

    in := struct {
        AdminUser     string `json:"admin_user"`
        AdminPassword string `json:"admin_password"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "ERROR: Fail to decode json for updating intance", http.StatusBadRequest)
        return
    }

    instance.AdminUser = in.AdminUser
    instance.AdminPassword = in.AdminPassword
    instance.Update()

    json.NewEncoder(w).Encode(instance)
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
    instance := parseInstance(r)
    if instance == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    instance.Delete()
}

////////////////////////////////////////////////////////////////////////////////

func parseInstance(r *http.Request) *model.Instance {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        panic(err)
    }

    return model.GetInstance(id)
}
