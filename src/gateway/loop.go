package main

import (
    "fmt"
    "log"
    "time"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    "gateway/model"
)

func checkHealth(host string, port int, user, password, database string) bool {
    uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
    db, err := sql.Open("mysql", uri)
    if err != nil {
        return false
    }

    if db.Ping() != nil {
        return false
    }

    return true
}

func mysqlHealthLoop() {
    for _, s := range model.ListInstance(nil, nil, nil) {
        good := checkHealth(s.Host, s.Port, s.AdminUser, s.AdminPassword, "mysql")
        if good == s.Active {
            continue
        }

        log.Printf(
            "INFO: mysql server `%s:%s` active status change from `%v` to `%v`\n",
            s.Host, s.Port, s.Active, good,
        )
    }

    for _, mi := range model.ListMapping(nil, nil, nil) {
        good := checkHealth(mi.MysqlHost, mi.MysqlPort, mi.MysqlUser, mi.MysqlPassword, mi.MysqlDb)
        if good == mi.MysqlActive {
            continue
        }

        log.Printf(
            "INFO: mysql instance `%s:%s-%s` active status change from `%v` to `%v`\n",
            mi.MysqlHost, mi.MysqlPort, mi.MysqlDb, mi.MysqlActive, good,
        )
    }

    time.Sleep(2 * time.Second)
}
