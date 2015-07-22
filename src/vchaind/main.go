package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "vchaind/model"
)

var conf *Config

type Database struct {
    Db       string `json:"db"`
    User     string `json:"user"`
    Password string `json:"password"`
    Host     string `json:"host"`
    Port     string `json:"port"`
}

type Config struct {
    Port     string   `json:"port"`
    Log      string   `json:"log"`
    Database Database `json:"database"`
}

func parseConfig(path string) {
    bytes, e := ioutil.ReadFile(path)
    if e != nil {
        log.Fatalf("Failed to read config file, %s\n", e.Error())
    }

    conf = new(Config)
    e = json.Unmarshal(bytes, conf)
    if e != nil {
        log.Fatalf("Failed to parse config file, %s\n", e.Error())
    }
}

func initialize() {
    if conf.Port == "" {
        conf.Port = "80"
    }

    if conf.Log != "" {
        f, e := os.OpenFile(conf.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if e != nil {
            log.Fatalf("Failed to open log file, %s\n", e.Error())
        }

        log.Printf("Redirecting log to file %s\n", conf.Log)
        log.SetOutput(f)
    }

    if err := model.ConnectDatabase(
        conf.Database.Host,
        conf.Database.Port,
        conf.Database.User,
        conf.Database.Password,
        conf.Database.Db,
    ); err != nil {
        log.Fatalf("Failed to connect database, %s\n", err.Error())
    }
}

func main() {
    path := flag.String("c", "/vchain/config/dispatcher.json", "config file")
    flag.Parse()

    parseConfig(*path)
    initialize()

    log.Printf("Starting vchaind on 0.0.0.0:%s ...\n", conf.Port)

    router := NewRouter()
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), router))
}
