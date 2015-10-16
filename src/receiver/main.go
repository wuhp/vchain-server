package main

import (
    "os"
    "fmt"
    "flag"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "discover"
)

var conf *Config

type Config struct {
    Port     string `json:"port"`
    Log      string `json:"log"`
    Consumer string `json:"consumer"`
    Vchaind  string `json:"vchaind"`
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

    if conf.Port == "" {
        conf.Port = "80"
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

    discover.InitConsumer(conf.Consumer)
    discover.InitVchaind(conf.Vchaind)
}

func main() {
    path := flag.String("c", "/vchain/server/config/consumer.json", "config file")
    flag.Parse()

    parseConfig(*path)
    initialize()

    log.Printf("Starting consumer on 0.0.0.0:%s ...\n", conf.Port)

    router := NewRouter()
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), router))
}
