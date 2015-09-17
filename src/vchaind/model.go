package main

import (
    "fmt"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDatabase(host, port, user, password, database string) (err error) {
    uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

    if db, err = sql.Open("mysql", uri); err != nil {
        return
    }

    err = db.Ping()
    return
}
