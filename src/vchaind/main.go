package main

import (
    "os"
    "fmt"
    "flag"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
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

    if err := ConnectDatabase(
        conf.Database.Host,
        conf.Database.Port,
        conf.Database.User,
        conf.Database.Password,
        conf.Database.Db,
    ); err != nil {
        log.Fatalf("ERROR: Failed to connect to database `%v`, with err `%s`\n", conf.Database, err.Error())
    }
}

func main() {
    path := flag.String("c", "/vchaind/config/account.json", "config file")
    flag.Parse()

    parseConfig(*path)
    initialize()

    log.Printf("Starting account on 0.0.0.0:%s ...\n", conf.Port)

    router := NewRouter()
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), router))
}
