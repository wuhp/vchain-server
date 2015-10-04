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

    defer db.Close()

    if db.Ping() != nil {
        return false
    }

    return true
}

func mysqlHealthLoop() {
    for {
        for _, s := range model.ListInstance(nil, nil, nil) {
            good := checkHealth(s.Host, s.Port, s.AdminUser, s.AdminPassword, "mysql")
            if good == s.Active {
                continue
            }

            s.Active = good
            s.Update()

            status := "active"
            if !good {
                status = "inactive"
            }

            log.Printf(
                "INFO: mysql server `%s:%d` become %s\n",
                s.Host, s.Port, status,
            )
        }

        for _, mi := range model.ListMapping(nil, nil, nil) {
            good := checkHealth(mi.MysqlHost, mi.MysqlPort, mi.MysqlUser, mi.MysqlPassword, mi.MysqlDb)
            if good == mi.MysqlActive {
                continue
            }

            mi.MysqlActive = good
            mi.Update()

            status := "active"
            if !good {
                status = "inactive"
            }

            log.Printf(
                "INFO: mysql instance `%s:%d@%s` become %s\n",
                mi.MysqlHost, mi.MysqlPort, mi.MysqlDb, status,
            )
        }

        time.Sleep(2 * time.Second)
    }
}
