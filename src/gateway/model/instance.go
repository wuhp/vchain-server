package model

import (
    "time"
)

type Instance struct {
    Id            int64  `json:"id"`
    Host          string `json:"host"`
    Port          int    `json:"port"`
    AdminUser     string `json:"admin_user"`
    AdminPassword string `json:"admin_password"`
    Active        bool   `json:"active"`
    CreateTs      int64  `json:"create_ts"`
}

////////////////////////////////////////////////////////////////////////////////

func ListInstance(cs []*Condition, o *Order, p *Paging) []*Instance {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, host, port, admin_user, admin_password,
            active, create_ts
        FROM
            instances
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Instance, 0)
    for rows.Next() {
        r := new(Instance)
        if err := rows.Scan(
            &r.Id, &r.Host, &r.Port, &r.AdminUser, &r.AdminPassword,
            &r.Active, &r.CreateTs,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func GetInstance(id int64) *Instance {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListInstance(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Instance) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO instances(
            host, port, admin_user, admin_password, active,
            create_ts
        )
        VALUES(?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    r.CreateTs = time.Now().UTC().Unix()

    result, err := stmt.Exec(
        r.Host, r.Port, r.AdminUser, r.AdminPassword, r.Active,
        r.CreateTs,
    )
    if err != nil {
        panic(err)
    }

    r.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (r *Instance) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            instances
        SET
            host = ?,
            port = ?,
            admin_user = ?,
            admin_password = ?,
            active = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.Host,
        r.Port,
        r.AdminUser,
        r.AdminPassword,
        r.Active,
        r.Id,
    ); err != nil {
        panic(err)
    }
}

func (r *Instance) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            instances
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Id); err != nil {
        panic(err)
    }
}
