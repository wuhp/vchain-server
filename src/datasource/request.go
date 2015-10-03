package datasource

import (
    "database/sql"
)

type Request struct {
    Uuid       string `json:"uuid"`
    ParentUuid string `json:"parent_uuid"`
    Service    string `json:"service"`
    Category   string `json:"category"`
    Sync       bool   `json:"sync"`
    BeginTs    int64  `json:"begin_ts"`
    EndTs      int64  `json:"end_ts"`
    GroupUuid  string `json:"-"`
    CreateTs   int64  `json:"-"`
}

func ListRequest(db *sql.DB, cs []*Condition, o *Order, p *Paging) []*Request {
    where, vs := generateWhereSql(cs)
    order := generateOrderSql(o)
    limit := generateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            uuid, parent_uuid, service, category, sync,
            begin_ts, end_ts, group_uuid, create_ts
        FROM
            request
        `+where+order+limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Request, 0)
    for rows.Next() {
        r := new(Request)
        if err := rows.Scan(
            &r.Uuid, &r.ParentUuid, &r.Service, &r.Category, &r.Sync,
            &r.BeginTs, &r.EndTs, &r.GroupUuid, &r.CreateTs,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func GetRequest(db *sql.DB, uuid string) *Request {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("uuid", "=", uuid))

    l := ListRequest(db, conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Request) Save(db *sql.DB) {
    stmt, err := db.Prepare(`
        INSERT INTO request(
            uuid, parent_uuid, service, category, sync,
            begin_ts, end_ts, group_uuid, create_ts
        )
        VALUES(
            ?, ?, ?, ?, ?,
            ?, ?, ?, ?
        )
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.Uuid, r.ParentUuid, r.Service, r.Category, r.Sync,
        r.BeginTs, r.EndTs, r.GroupUuid, r.CreateTs,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Update(db *sql.DB) {
    stmt, err := db.Prepare(`
        UPDATE
            request
        SET
            group_uuid = ?
        WHERE
            uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.GroupUuid,
        r.Uuid,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Delete(db *sql.DB) {
    stmt, err := db.Prepare(`
        DELETE FROM request
        WHERE uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Uuid); err != nil {
        panic(err)
    }
}
