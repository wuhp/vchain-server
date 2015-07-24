package model

import (
    "fmt"
    "strings"
)

type Service struct {
    Uuid       string `json:"uuid"`
    Category   string `json:"category"`
    Instance   string `json:"instance"`
    Hostname   string `json:"hostname"`
    StartTs    int64  `json:"start_ts"`
    StopTs     int64  `json:"stop_ts"`
}

func ListService(cs []*Condition, o *Order, p *Paging) []*Service {
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
            category,
            instance,
            hostname,
            start_ts,
            stop_ts
        FROM
            service
        `+where+order+limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Service, 0)
    for rows.Next() {
        s := new(Service)
        if err := rows.Scan(
            &s.Uuid,
            &s.Category,
            &s.Instance,
            &s.Hostname,
            &s.StartTs,
            &s.StopTs,
        ); err != nil {
            panic(err)
        }

        l = append(l, s)
    }

    return l
}

func (s *Service) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO
        service(
            uuid,
            category,
            instance,
            hostname,
            start_ts,
            stop_ts
        )
        VALUES(?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        s.Uuid,
        s.Category,
        s.Instance,
        s.Hostname,
        s.StartTs,
        s.StopTs,
    ); err != nil {
        panic(err)
    }
}

func (s *Service) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            service
        SET
            category = ?,
            instance = ?,
            hostname = ?,
            start_ts = ?,
            stop_ts = ?
        WHERE
            uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        s.Category,
        s.Instance,
        s.Hostname,
        s.StartTs,
        s.StopTs,
        s.Uuid,
    ); err != nil {
        panic(err)
    }
}

func (s *Service) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            service
        WHERE
            uuid=?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(s.Uuid); err != nil {
        panic(err)
    }
}
