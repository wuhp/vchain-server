package handler

import (
    "fmt"
    "log"
    "time"
    "strconv"
    "io/ioutil"
    "net/http"
    "database/sql"
    "encoding/json"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"

    "discover"
    "datasource"
)

func fetchDbConnection(gateway string, pid int64) *sql.DB {
    method := "GET"
    url := fmt.Sprintf("http://%s/api/v1/mysql/%d", gateway, pid)
    client := &http.Client{}

    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to create request, %s %s, with err %v\n", method, url, err))
    }

    res, err := client.Do(req)
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to send request, %s %s, with err %v\n", method, url, err))
    }
    if res.StatusCode == http.StatusNotFound {
        return nil
    }
    if res.StatusCode != http.StatusOK {
        panic(fmt.Sprintf("ERROR: Bad response status, %s %s, with err %s\n", method, url, res.Status))
    }

    defer res.Body.Close()
    out, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to read response body, %s %s, with err %v\n", method, url, err))
    }

    log.Printf("INFO: %s %s %s\n", method, url, res.Status)

    mapping := struct {
        ProjectId     int64  `json:"project_id"`
        MysqlHost     string `json:"mysql_host"`
        MysqlPort     int    `json:"mysql_port"`
        MysqlUser     string `json:"mysql_user"`
        MysqlPassword string `json:"mysql_password"`
        MysqlDb       string `json:"mysql_db"`
        MysqlActive   bool   `json:"mysql_active"`
        CreateTs      int64  `json:"create_ts"`
    } {}

    if err := json.Unmarshal(out, &mapping); err != nil {
        panic(fmt.Sprintf("ERROR: Fail to parse response body, %s %s, with err %v\n", method, url, err))
    }

    // init db object
    uri := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s",
        mapping.MysqlUser,
        mapping.MysqlPassword,
        mapping.MysqlHost,
        mapping.MysqlPort,
        mapping.MysqlDb,
    )

    db, err := sql.Open("mysql", uri)
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to connect to mysql, %s, with err %v\n", uri, err))
    }

    return db
}

func getDb(r *http.Request) *sql.DB {
    vars := mux.Vars(r)
    pid, _ := strconv.ParseInt(vars["pid"], 10, 64)
    gateway, err := discover.DiscoverGateway()
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to discover gateway, with err %v\n", err))
    }

    return fetchDbConnection(gateway, pid)
}

func getId(r *http.Request) string {
    vars := mux.Vars(r)
    return vars["id"]
}

func getUuid(r *http.Request) string {
    vars := mux.Vars(r)
    return vars["uuid"]
}

func getRequestType(r *http.Request) *datasource.RequestType {
    vars := mux.Vars(r)

    return &datasource.RequestType{
        Service: vars["service"],
        Category: vars["category"],
    }
}

func getTimeRange(r *http.Request) *datasource.TimeRange {
    from := r.URL.Query().Get("ts_from")
    to := r.URL.Query().Get("ts_to")

    tr := new(datasource.TimeRange)
    if len(to) == 0 {
        tr.End = time.Now().UTC().UnixNano()
    } else {
        tr.End, _ = strconv.ParseInt(to, 10, 64)
    }

    if len(from) == 0 {
        tr.Begin = time.Now().UTC().UnixNano() - 24 * time.Hour.Nanoseconds()
    } else {
        tr.Begin, _ = strconv.ParseInt(from, 10, 64)
    }

    return tr
}
