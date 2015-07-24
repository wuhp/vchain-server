package model

import (
    "fmt"
    "strings"
)

type Request struct {
    Uuid          string `json:"uuid"`
    ServiceUuid   string `json:"service_uuid"`
    ParentUuid    string `json:"parent_uuid"`
    Category      string `json:"category"`
    BeginTs       int64  `json:"begin_ts"`
    EndTs         int64  `json:"end_ts"`
    BeginMetadata string `json:"begin_metadata"`
    EndMetadata   string `json:"end_metadata"`
}

func ListRequest(cs []*Condition, o *Order, p *Paging) []*Request {
    ks := make([]string, 0)
    vs := make([]interface{}, 0)
    for _, c := range cs {
        ks = append(ks, fmt.Sprintf("%s%s?", c.Key, c.Op))
        vs = append(vs, c.Value)
    }

    where := ""
    if cs != nil {
        where = "WHERE " + strings.Join(ks, " and ")
    }

    order := ""
    if o != nil {
        order = fmt.Sprintf(" ORDER BY %s %s", o.Columns, o.Sequence)
    }

    limit := ""
    if p != nil {
        limit = fmt.Sprintf(" LIMIT %d, %d", p.Offset, p.Size)
    }

    rows, err := db.Query(`
        SELECT
            uuid,
            service_uuid,
            parent_uuid,
            category,
            begin_ts,
            end_ts,
            begin_metadata,
            end_metadata
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
            &r.Uuid,
            &r.ServiceUuid,
            &r.ParentUuid,
            &r.Category,
            &r.BeginTs,
            &r.EndTs,
            &r.BeginMetadata,
            &r.EndMetadata,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func (r *Request) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO
        request(
            uuid,
            service_uuid,
            parent_uuid,
            category,
            begin_ts,
            end_ts,
            begin_metadata,
            end_metadata
        )
        VALUES(?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.Uuid,
        r.ServiceUuid,
        r.ParentUuid,
        r.Category,
        r.BeginTs,
        r.EndTs,
        r.BeginMetadata,
        r.EndMetadata,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            request
        SET
            service_uuid = ?,
            parent_uuid = ?,
            category = ?,
            begin_ts = ?,
            end_ts = ?,
            begin_metadata = ?,
            end_metadata = ?
        WHERE
            uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.ServiceUuid,
        r.ParentUuid,
        r.Category,
        r.BeginTs,
        r.EndTs,
        r.BeginMetadata,
        r.EndMetadata,
        r.Uuid,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            request
        WHERE
            uuid=?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Uuid); err != nil {
        panic(err)
    }
}
