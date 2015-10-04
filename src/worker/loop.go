package main

import (
    "log"
    "fmt"
    "time"
    "io/ioutil"
    "net/http"
    "encoding/json"

    _ "github.com/go-sql-driver/mysql"

    "discover"
    gm "gateway/model"
)

func mainLoop() {
    for {
        projects := getActiveProjects()
        for _, p := range projects {
            log.Printf("INFO: Start to process project %d\n", p.id)
            p.connect()
            p.process()
            p.disconnect()
        }
        time.Sleep(1 * time.Second)
    }
}

func getMappings(gateway string) []*gm.Mapping{
    method := "GET"
    url := fmt.Sprintf("http://%s/api/v1/mysql", gateway)
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

    mappings := make([]*gm.Mapping, 0)
    if err := json.Unmarshal(out, &mappings); err != nil {
        panic(fmt.Sprintf("ERROR: Fail to parse response body, %s %s, with err %v\n", method, url, err))
    }

    return mappings
}

func getActiveProjects() []*Project {
    gateway, err := discover.DiscoverGateway()
    if err != nil {
        panic(fmt.Sprintf("ERROR: Fail to discover gateway, with err %v\n", err))
    }

    mappings := getMappings(gateway)
    projects := make([]*Project, 0)
    for _, mapping := range mappings {
        if !mapping.MysqlActive {
            log.Printf("INFO: Ignore to process data for project %d, mysql is not available\n", mapping.ProjectId)
            continue
        }

        p := new(Project)
        p.id = mapping.ProjectId
        p.host = mapping.MysqlHost
        p.port = mapping.MysqlPort
        p.user = mapping.MysqlUser
        p.password = mapping.MysqlPassword
        p.database = mapping.MysqlDb
        projects = append(projects, p)
    }

    return projects
}
