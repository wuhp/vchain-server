package datasource

import (
    "database/sql"
)

func GetServices(db *sql.DB, tr *TimeRange) []string {
    sql := `
        SELECT DISTINCT(service)
        FROM request
        WHERE begin_ts > ? AND begin_ts < ? AND group_uuid != ""
    `

    rows, err := db.Query(sql, tr.Begin, tr.End)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]string, 0)
    for rows.Next() {
        var s string
        if err := rows.Scan(&s); err != nil {
            panic(err)
        }

        l = append(l, s)
    }

    return l
}

func GetServiceChain(db *sql.DB, tr *TimeRange) []*Pair {
    sql := `
        SELECT r1.service, r2.service
        FROM request AS r1, request AS r2
        WHERE r1.uuid = r2.parent_uuid AND r1.group_uuid != "" AND r1.begin_ts > ? AND r1.begin_ts < ?
        GROUP BY r1.service, r2.service
    `

    rows, err := db.Query(sql, tr.Begin, tr.End)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Pair, 0)
    for rows.Next() {
        p := new(Pair)
        if err := rows.Scan(&p.From, &p.To); err != nil {
            panic(err)
        }

        l = append(l, p)
    }

    return l
}

func GetRequestTypes(db *sql.DB, tr *TimeRange) []*RequestType {
    sql := `
        SELECT service, category
        FROM request
        WHERE begin_ts > ? AND begin_ts < ?
        GROUP BY service, category
    `

    rows, err := db.Query(sql, tr.Begin, tr.End)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*RequestType, 0)
    for rows.Next() {
        rt := new(RequestType)
        if err := rows.Scan(&rt.Service, &rt.Category); err != nil {
            panic(err)
        }

        l = append(l, rt)
    }

    return l
}

func QueryInvokeChain(db *sql.DB, tr *TimeRange, rt *RequestType) []*InvokeChain {
    vs := make([]interface{}, 0)
    sql := `
        SELECT ic.id, ic.header, ic.request_types, ic.parents_index
        FROM request AS r, request_group AS rg, invoke_chain AS ic
        WHERE ic.id = rg.invoke_chain_id AND rg.uuid = r.uuid AND r.begin_ts > ? AND r.begin_ts < ?
    `
    vs = append(vs, tr.Begin)
    vs = append(vs, tr.End)

    if rt != nil {
        sql = sql + ` AND header = ?`
        vs = append(vs, RequestType2string(rt))
    }

    rows, err := db.Query(sql, vs...)

    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*InvokeChain, 0)
    for rows.Next() {
        ic := new(InvokeChain)
        var h string
        var rts string
        var pis string
        if err := rows.Scan(&ic.Id, &h, &rts, &pis); err != nil {
            panic(err)
        }

        ic.Header = *String2requestType(h)
        ic.RequestTypes = String2requestTypes(rts)
        ic.ParentsIndex = String2ints(pis)
        l = append(l, ic)
    }

    return l
}

func FindRequestsByInvokeChain(db *sql.DB, invokeChainId int64) []*Request {
    sql := `
        SELECT
            r.uuid, r.parent_uuid, r.service, r.category, r.sync,
            r.begin_ts, r.end_ts, r.group_uuid, r.create_ts
        FROM
            request AS r, request_group AS rg, invoke_chain AS ic
        WHERE
            ic.id = ? AND ic.id = rg.invoke_chain_id AND rg.uuid = r.uuid
    `

    rows, err := db.Query(sql, invokeChainId)
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
