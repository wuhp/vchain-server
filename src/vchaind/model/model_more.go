package model

import (
    "encoding/json"
)

func GetServices() []string {
    sql := `
        SELECT DISTINCT(service)
        FROM request
    `

    rows, err := db.Query(sql)
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

func GetServiceChain() []*Pair {
    sql := `
        SELECT r1.service, r2.service
        FROM request AS r1, request AS r2
        WHERE r1.uuid = r2.parent_uuid AND r1.group_uuid != ""
        GROUP BY r1.service, r2.service
    `

    rows, err := db.Query(sql)
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

func GetRequestTypes() []*RequestType {
    sql := `
        SELECT service, category
        FROM request
        GROUP BY service, category
    `   
        
    rows, err := db.Query(sql)
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

func FindRequestsByInvokeChain(invokeChainId int64) []*Request {
    sql := `
        SELECT
            r.uuid, r.parent_uuid, r.service, r.category, r.sync_option, r.begin_ts, r.end_ts,
            r.begin_metadata, r.end_metadata, r.group_uuid, r.create_ts, r.update_ts
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
        var bm string
        var em string
        if err := rows.Scan(
            &r.Uuid, &r.ParentUuid, &r.Service, &r.Category, &r.SyncOption, &r.BeginTs, &r.EndTs,
            &bm, &em, &r.GroupUuid, &r.CreateTs, &r.UpdateTs,
        ); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(bm), &r.BeginMetadata); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(em), &r.EndMetadata); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}
