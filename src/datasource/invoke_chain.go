package datasource

import (
    "database/sql"
)

type InvokeChain struct {
    Id           int64          `json:"id"`
    Header       RequestType    `json:"-"`
    RequestTypes []*RequestType `json:"request_types"`
    ParentsIndex []int          `json:"parents_index"`
}

func ListInvokeChain(db *sql.DB, cs []*Condition, o *Order, p *Paging) []*InvokeChain {
    where, vs := generateWhereSql(cs)
    order := generateOrderSql(o)
    limit := generateLimitSql(p)

    rows, err := db.Query(`
        SELECT id, header, request_types, parents_index
        FROM invoke_chain
    `+where+order+limit, vs...)

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

func GetInvokeChain(db *sql.DB, id int64) *InvokeChain {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListInvokeChain(db, conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetInvokeChainByValues(db *sql.DB, types []*RequestType, index []int) *InvokeChain {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("request_types", "=", RequestTypes2string(types)))
    conditions = append(conditions, NewCondition("parents_index", "=", Ints2string(index)))

    l := ListInvokeChain(db, conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (ic *InvokeChain) Save(db *sql.DB) {
    stmt, err := db.Prepare(`
        INSERT INTO invoke_chain(header, request_types, parents_index)
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    h := RequestType2string(&ic.Header)
    rts := RequestTypes2string(ic.RequestTypes)
    pis := Ints2string(ic.ParentsIndex)

    result, err := stmt.Exec(h, rts, pis)
    if err != nil {
        panic(err)
    }

    ic.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (ic *InvokeChain) Delete(db *sql.DB) {
    stmt, err := db.Prepare(`
        DELETE FROM invoke_chain
        WHERE id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(ic.Id); err != nil {
        panic(err)
    }
}


