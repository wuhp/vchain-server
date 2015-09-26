package model

import (
    "time"
)

type Mapping struct {
    ProjectId     int64  `json:"project_id"`
    MysqlHost     string `json:"mysql_host"`
    MysqlPort     int    `json:"mysql_port"`
    MysqlUser     string `json:"mysql_user"`
    MysqlPassword string `json:"mysql_password"`
    MysqlDb       string `json:"mysql_db"`
    MysqlActive   bool   `json:"mysql_active"`
    CreateTs      int64  `json:"create_ts"`
}

////////////////////////////////////////////////////////////////////////////////

func ListMapping(cs []*Condition, o *Order, p *Paging) []*Mapping {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            project_id, mysql_host, mysql_port, mysql_user, mysql_password,
            mysql_db, mysql_active, create_ts
        FROM
            mappings
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Mapping, 0)
    for rows.Next() {
        r := new(Mapping)
        if err := rows.Scan(
            &r.ProjectId, &r.MysqlHost, &r.MysqlPort, &r.MysqlUser, &r.MysqlPassword,
            &r.MysqlDb, &r.MysqlActive, &r.CreateTs,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func GetMapping(pid int64) *Mapping {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("project_id", "=", pid))

    l := ListMapping(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Mapping) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO mappings(
            project_id, mysql_host, mysql_port, mysql_user, mysql_password,
            mysql_db, mysql_active, create_ts
        )
        VALUES(?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    r.CreateTs = time.Now().UTC().Unix()

    if _, err := stmt.Exec(
        r.ProjectId, r.MysqlHost, r.MysqlPort, r.MysqlUser, r.MysqlPassword,
        r.MysqlDb, r.MysqlActive, r.CreateTs,
    ); err != nil {
        panic(err)
    }
}

func (r *Mapping) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            mappings
        SET
            mysql_host = ?,
            mysql_port = ?,
            mysql_user = ?,
            mysql_password = ?,
            mysql_db = ?,
            mysql_active = ?
        WHERE
            project_id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.MysqlHost,
        r.MysqlPort,
        r.MysqlUser,
        r.MysqlPassword,
        r.MysqlDb,
        r.MysqlActive,
        r.ProjectId,
    ); err != nil {
        panic(err)
    }
}

func (r *Mapping) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            mappings
        WHERE
            project_id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.ProjectId); err != nil {
        panic(err)
    }
}
