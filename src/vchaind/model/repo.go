package model

import (
    "time"
)

type Repo struct {
    Id       int64  `json:"id"`
    UserId   int64  `json:"user_id"`
    Name     string `json:"name"`
    Hash     string `json:"hash"`
    CreateTs int64  `json:"create_ts"`
}

////////////////////////////////////////////////////////////////////////////////

func ListRepo(cs []*Condition, o *Order, p *Paging) []*Repo {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, user_id, name, hash, create_ts
        FROM
            repos
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Repo, 0)
    for rows.Next() {
        r := new(Repo)
        if err := rows.Scan(
            &r.Id, &r.UserId, &r.Name, &r.Hash, &r.CreateTs,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func GetRepo(id int64) *Repo {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListRepo(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Repo) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO repos(
            user_id, name, hash, create_ts
        )
        VALUES(?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    r.Hash = generateHash()
    r.CreateTs = time.Now().UTC().Unix()

    result, err := stmt.Exec(
        r.UserId, r.Name, r.Hash, r.CreateTs,
    )
    if err != nil {
        panic(err)
    }

    r.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (r *Repo) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            repos
        SET
            name = ?,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.Name,
        r.Id,
    ); err != nil {
        panic(err)
    }
}

func (r *Repo) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            repos
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

////////////////////////////////////////////////////////////////////////////////

func (r *Repo) ResetHash() {
    stmt, err := db.Prepare(`
        UPDATE
            repos
        SET
            hash = ?,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    r.Hash = generateHash()

    if _, err := stmt.Exec(
        r.Hash,
        r.Id,
    ); err != nil {
        panic(err)
    }
}
