package main

import (
    "os"
    "flag"
    "log"
    "io/ioutil"
    "encoding/json"

    "discover"
)

var conf *Config

type Config struct {
    Log      string   `json:"log"`
    Gateway  string   `json:"gateway"`
}

func parseConfig(path string) {
    bytes, e := ioutil.ReadFile(path)
    if e != nil {
        log.Fatalf("ERROR: Failed to read config file `%s`, with err `%s`\n", path, e.Error())
    }

    conf = new(Config)
    e = json.Unmarshal(bytes, conf)
    if e != nil {
        log.Fatalf("ERROR: Failed to parse config file `%s`, with err `%s`\n", path, e.Error())
    }
}

func initialize() {
    if conf.Log != "" {
        f, e := os.OpenFile(conf.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if e != nil {
            log.Fatalf("ERROR: Failed to open log file `%s`, with err `%s`\n", conf.Log, e.Error())
        }

        log.SetOutput(f)
    }

    discover.InitGateway(conf.Gateway)
}

func main() {
    path := flag.String("c", "/vchain/server/config/worker.json", "config file")
    flag.Parse()

    parseConfig(*path)
    initialize()

    mainLoop()
}
