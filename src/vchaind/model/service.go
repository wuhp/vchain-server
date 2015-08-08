package model

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
        SELECT DISTINCT(r1.service, r2.service)
        FROM request AS r1, request AS r2
        WHERE r1.uuid = r2.parent_uuid AND r1.group_uuid != ""
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
