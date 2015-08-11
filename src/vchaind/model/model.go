package model

import (
    "fmt"
    "encoding/json"
)

/////////////////////////////////////////////////////////////////////////////////

type Pair struct {
    From string `json:"from"`
    To   string `json:"to"`
}

func (p *Pair) String() string {
    return fmt.Sprintf("(%s,%s)", p.From, p.To)
}

type RequestType struct {
    Service  string `json:"service"`
    Category string `json:"category"`
}

func (rt *RequestType) String() string {
    return fmt.Sprintf("(%s,%s)", rt.Service, rt.Category)
}

///////////////////////////////////////////////////////////////////////////////

type Request struct {
    Uuid          string            `json:"uuid"`
    ParentUuid    string            `json:"parent_uuid"`
    Service       string            `json:"service"`
    Category      string            `json:"category"`
    SyncOption    string            `json:"sync_option"`
    BeginTs       int64             `json:"begin_ts"`
    EndTs         int64             `json:"end_ts"`
    BeginMetadata map[string]string `json:"begin_metadata"`
    EndMetadata   map[string]string `json:"end_metadata"`
    GroupUuid     string            `json:"-"`
    CreateTs      int64             `json:"-"`
    UpdateTs      int64             `json:"-"`
}

type RequestGroup struct {
    Uuid          string   `json:"uuid"`
    RequestUuids  []string `json:"request_uuids"`
    ParentsIndex  []int    `json:"parents_index"`
    InvokeChainId int64    `json:"-"`
}

type InvokeChain struct {
    Id           int64          `json:"id"`
    Header       RequestType    `json:"-"`
    RequestTypes []*RequestType `json:"request_types"`
    ParentsIndex []int          `json:"parents_index"`
}

////////////////////////////////////////////////////////////////////////////////

func ListRequest(cs []*Condition, o *Order, p *Paging) []*Request {
    where, vs := generateWhereSql(cs)
    order := generateOrderSql(o)
    limit := generateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            uuid, parent_uuid, service, category, sync_option, begin_ts, end_ts,
            begin_metadata, end_metadata, group_uuid, create_ts, update_ts
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

func GetRequest(uuid string) *Request {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("uuid", "=", uuid))

    l := ListRequest(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Request) Exist() bool {
    return nil != GetRequest(r.Uuid)
}

func (r *Request) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO request(
            uuid, parent_uuid, service, category, sync_option, begin_ts, end_ts,
            begin_metadata, end_metadata, group_uuid, create_ts, update_ts
        )
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    bm, _ := json.Marshal(&r.BeginMetadata)
    em, _ := json.Marshal(&r.EndMetadata)

    if _, err := stmt.Exec(
        r.Uuid, r.ParentUuid, r.Service, r.Category, r.SyncOption, r.BeginTs, r.EndTs,
        string(bm), string(em), r.GroupUuid, r.CreateTs, r.UpdateTs,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            request
        SET
            parent_uuid = ?,
            service = ?,
            category = ?,
            sync_option = ?,
            begin_ts = ?,
            end_ts = ?,
            begin_metadata = ?,
            end_metadata = ?,
            group_uuid = ?,
            create_ts = ?,
            update_ts = ?
        WHERE
            uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    bm, _ := json.Marshal(&r.BeginMetadata)
    em, _ := json.Marshal(&r.EndMetadata)

    if _, err := stmt.Exec(
        r.ParentUuid,
        r.Service,
        r.Category,
        r.SyncOption,
        r.BeginTs,
        r.EndTs,
        string(bm),
        string(em),
        r.GroupUuid,
        r.CreateTs,
        r.UpdateTs,
        r.Uuid,
    ); err != nil {
        panic(err)
    }
}

func (r *Request) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM request
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

////////////////////////////////////////////////////////////////////////////////

func ListRequestGroup(cs []*Condition, o *Order, p *Paging) []*RequestGroup {
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

func GetRequestGroup(uuid string) *RequestGroup {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("uuid", "=", uuid))

    l := ListRequestGroup(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *RequestGroup) Exist() bool {
    return nil != GetRequestGroup(r.Uuid)
}

func (r *RequestGroup) DetailRequests() []*Request {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("group_uuid", "=", r.Uuid))
    group := ListRequest(conditions, nil, nil)

    result := make([]*Request, 0)
    for _, uuid := range r.RequestUuids {
        i := FindRequestByUuid(group, uuid)
        if i < 0 {
            // TBD
        }
        result = append(result, group[i])
    }

    return result
}

func (r *RequestGroup) Save() {
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

func (r *RequestGroup) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            request_group
        SET
            request_uuids = ?,
            parents_index = ?,
            invoke_chain_id = ?
        WHERE
            uuid = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    rus := Strings2string(r.RequestUuids)
    pis := Ints2string(r.ParentsIndex)

    if _, err := stmt.Exec(rus, pis, r.InvokeChainId, r.Uuid); err != nil {
        panic(err)
    }
}

func (r *RequestGroup) Delete() {
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

////////////////////////////////////////////////////////////////////////////////

func ListInvokeChain(cs []*Condition, o *Order, p *Paging) []*InvokeChain {
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

func GetInvokeChain(id int64) *InvokeChain {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListInvokeChain(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetInvokeChainByValues(types []*RequestType, index []int) *InvokeChain {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("request_types", "=", RequestTypes2string(types)))
    conditions = append(conditions, NewCondition("parents_index", "=", Ints2string(index)))

    l := ListInvokeChain(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (ic *InvokeChain) Save() {
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

func (ic *InvokeChain) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            invoke_chain
        SET
            header = ?,
            request_types = ?,
            parents_index = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    h := RequestType2string(&ic.Header)
    rts := RequestTypes2string(ic.RequestTypes)
    pis := Ints2string(ic.ParentsIndex)

    if _, err := stmt.Exec(h, rts, pis, ic.Id); err != nil {
        panic(err)
    }
}

func (ic *InvokeChain) Delete() {
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


