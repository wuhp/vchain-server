package datasource

import (
    "database/sql"
)

type RequestGroup struct {
    Uuid          string   `json:"uuid"`
    RequestUuids  []string `json:"request_uuids"`
    ParentsIndex  []int    `json:"parents_index"`
    InvokeChainId int64    `json:"-"`
}

func ListRequestGroup(db *sql.DB, cs []*Condition, o *Order, p *Paging) []*RequestGroup {
    where, vs := generateWhereSql(cs)
    order := generateOrderSql(o)
    limit := generateLimitSql(p)

    rows, err := db.Query(`
        SELECT uuid, request_uuids, parents_index, invoke_chain_id
        FROM request_group
    `+where+order+limit, vs...)

    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*RequestGroup, 0)
    for rows.Next() {
        r := new(RequestGroup)
        var rus string
        var pis string
        if err := rows.Scan(
            &r.Uuid, &rus, &pis, &r.InvokeChainId,
        ); err != nil {
            panic(err)
        }

        r.RequestUuids = String2strings(rus)
        r.ParentsIndex = String2ints(pis)
        l = append(l, r)
    }

    return l
}

func GetRequestGroup(db *sql.DB, uuid string) *RequestGroup {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("uuid", "=", uuid))

    l := ListRequestGroup(db, conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *RequestGroup) Save(db *sql.DB) {
    stmt, err := db.Prepare(`
        INSERT INTO request_group(uuid, request_uuids, parents_index, invoke_chain_id)
        VALUES(?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    rus := Strings2string(r.RequestUuids)
    pis := Ints2string(r.ParentsIndex)

    if _, err := stmt.Exec(r.Uuid, rus, pis, r.InvokeChainId); err != nil {
        panic(err)
    }
}

func (r *RequestGroup) Delete(db *sql.DB) {
    stmt, err := db.Prepare(`
        DELETE FROM request_group
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

func (r *RequestGroup) DetailRequests(db *sql.DB) []*Request {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("group_uuid", "=", r.Uuid))
    group := ListRequest(db, conditions, nil, nil)

    result := make([]*Request, 0)
    for _, uuid := range r.RequestUuids {
        i := FindRequestByUuid(group, uuid)
        result = append(result, group[i])
    }

    return result
}
