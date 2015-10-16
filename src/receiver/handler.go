package main

import (
    "fmt"
    "log"
    "bytes"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "datasource"
    "discover"
)

type Data struct {
    Reqs  []*datasource.Request    `json:"requests"`
    Rlogs []*datasource.RequestLog `json:"request-logs"`
}

func ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func post(w http.ResponseWriter, r *http.Request) {
    pid, valid := verify(r.Header.Get("Authorization"))
    if !valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    defer r.Body.Close()
    data := new(Data)
    if err := json.NewDecoder(r.Body).Decode(data); err != nil {
        log.Printf("ERROR: can not decode payload, with err %s\n", err.Error())
        http.Error(w, "ERROR: can not decode payload", http.StatusBadRequest)
        return
    }

    // TBD
    // How to handle error here or even ignore it
    consume(pid, data)
}

func verify(key string) (int64, bool) {
    if key == "" {
        return 0, false
    }

    host, err := discover.DiscoverVchaind()
    if err != nil {
        log.Printf("ERROR: Fail to discover vchaind, with err %v\n", err)
        panic(err)
    }

    client := &http.Client{}
    req, err := http.NewRequest(
        "GET",
        fmt.Sprintf("http://%s/api/v1/verify/%s", host, key),
        nil,
    )
    if err != nil {
        log.Printf("ERROR: Fail to create http request, with err %v\n", err)
        panic(err)
    }

    req.Header.Set("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        log.Printf("ERROR: Fail to send http request, with err %v\n", err)
        panic(err)
    }

    log.Printf("INFO: Http response %s\n", res.Status)

    if res.StatusCode != http.StatusOK {
        panic("ERROR: Response status not 200")
    }

    defer res.Body.Close()

    project := struct {
        Ok  bool  `json:"ok"`
        Pid int64 `json:"pid"`
    } {}

    outbody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Printf("ERROR: Fail to read response body, with err %v\n", err)
        panic(err)
    }

    err = json.Unmarshal(outbody, &project)
    if err != nil {
        log.Printf("ERROR: Fail to decode response body, with err %v\n", err)
        panic(err)
    }

    return project.Pid, project.Ok
}

func consume(pid int64, data *Data) {
    host, err := discover.DiscoverConsumer()
    if err != nil {
        log.Printf("ERROR: Fail to discover consumer, with err %v\n", err)
        panic(err)
    }

    consumeRequest(host, pid, data.Reqs)
    consumeRequestLog(host, pid, data.Rlogs)
}

func consumeRequest(host string, pid int64, reqs []*datasource.Request) {
    client := &http.Client{}
    inbody, _ := json.Marshal(reqs)
    req, err := http.NewRequest(
        "POST",
        fmt.Sprintf("http://%s/api/v1/%d/requests", host, pid),
        bytes.NewBuffer(inbody),
    )
    if err != nil {
        log.Printf("ERROR: Fail to create http request, with err %v\n", err)
        panic(err)
    }

    req.Header.Set("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        log.Printf("ERROR: Fail to send http request, with err %v\n", err)
        panic(err)
    }

    log.Printf("INFO: Http response %s\n", res.Status)

    if res.StatusCode != http.StatusOK {
        panic("ERROR: Response status not 200")
    }
}

func consumeRequestLog(host string, pid int64, rlogs []*datasource.RequestLog) {
    client := &http.Client{}
    inbody, _ := json.Marshal(rlogs)
    req, err := http.NewRequest(
        "POST",
        fmt.Sprintf("http://%s/api/v1/%d/request-logs", host, pid),
        bytes.NewBuffer(inbody),
    )
    if err != nil {
        log.Printf("ERROR: Fail to create http request, with err %v\n", err)
        panic(err)
    }

    req.Header.Set("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        log.Printf("ERROR: Fail to send http request, with err %v\n", err)
        panic(err)
    }

    log.Printf("INFO: Http response %s\n", res.Status)

    if res.StatusCode != http.StatusOK {
        panic("ERROR: Response status not 200")
    }
}
