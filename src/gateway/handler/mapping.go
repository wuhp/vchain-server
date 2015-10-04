package handler

import (
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "gateway/model"
)

////////////////////////////////////////////////////////////////////////////////

func Provision(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        ProjectId int64 `json:"project_id"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "ERROR: Fail to decode json for mysql provison", http.StatusBadRequest)
        return
    }

    if model.GetMapping(in.ProjectId) != nil {
        w.WriteHeader(http.StatusConflict)
        return
    }

    mc := generateMysqlConnection(in.ProjectId)
    provisionMysqlInstance(mc)
    migrateMysqlInstance(mc)

    mapping := new(model.Mapping)
    mapping.ProjectId = in.ProjectId
    mapping.MysqlHost = mc.Host
    mapping.MysqlPort = mc.Port
    mapping.MysqlUser = mc.User
    mapping.MysqlPassword = mc.Password
    mapping.MysqlDb = mc.Database
    mapping.MysqlActive = true
    mapping.Save()

    json.NewEncoder(w).Encode(mapping)
}

func Unprovision(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        ProjectId int64 `json:"project_id"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "ERROR: Fail to decode json for mysql unprovison", http.StatusBadRequest)
        return
    }

    mapping := model.GetMapping(in.ProjectId)
    if mapping == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    mc := new(MysqlConnection)
    mc.Host = mapping.MysqlHost
    mc.Port = mapping.MysqlPort
    mc.User = mapping.MysqlUser
    mc.Password = mapping.MysqlPassword
    mc.Database = mapping.MysqlDb

    unprovisionMysqlInstance(mc)
    mapping.Delete()
}

func QueryMapping(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(model.ListMapping(nil, nil, nil))
}

func GetMapping(w http.ResponseWriter, r *http.Request) {
    mapping := parseMapping(r)
    if mapping == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(mapping)
}

////////////////////////////////////////////////////////////////////////////////

func parseMapping(r *http.Request) *model.Mapping {
    vars := mux.Vars(r)
    pid, err := strconv.ParseInt(vars["pid"], 10, 64)
    if err != nil {
        panic(err)
    }

    return model.GetMapping(pid)
}
