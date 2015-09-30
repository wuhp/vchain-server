package datasource

import (
    "database/sql"
)

type RequestLog struct {
    Uuid      string `json:"uuid"`
    Timestamp int64  `json:"timestamp"`
    Log       string `json:"log"`
}

func ListRequestLog(db *sql.DB, cs []*Condition, o *Order, p *Paging) []*RequestLog {
    where, vs := generateWhereSql(cs)
    order := generateOrderSql(o)
    limit := generateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            uuid, timestamp, log
        FROM
            request_log
        `+where+order+limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*RequestLog, 0)
    for rows.Next() {
        r := new(RequestLog)
        if err := rows.Scan(
            &r.Uuid, &r.Timestamp, &r.Log,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func DeleteRequestLog(db *sql.DB, cs []*Condition) {
    where, vs := generateWhereSql(cs)

    stmt, err := db.Prepare(`DELETE FROM request_log` + where)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(vs); err != nil {
        panic(err)
    }
}

func (r *RequestLog) Save(db *sql.DB) {
    stmt, err := db.Prepare(`
        INSERT INTO request_log(uuid, timestamp, log)
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Uuid, r.Timestamp, r.Log); err != nil {
        panic(err)
    }
}
