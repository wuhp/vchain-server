package main

import (
    "log"
    "fmt"
    "time"
    "database/sql"

    "datasource"
)

const (
    ProcessWaiting = 10
    ProcessTimeout = 10 * 60
)

type Project struct {
    id       int64
    user     string
    password string
    host     string
    port     int
    database string
    db       *sql.DB
}

func (p *Project) connect() {
    uri := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s",
        p.user,
        p.password,
        p.host,
        p.port,
        p.database,
    )

    db, err := sql.Open("mysql", uri)
    if err != nil {
        log.Printf("ERROR: Fail to connect to mysql %s for project %d, with err %v\n", uri, p.id, err)
        panic(err)
    }

    p.db = db
}

func (p *Project) disconnect() {
    if p.db != nil {
        p.db.Close()
    }
}

func (p *Project) process() {
    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("group_uuid", "=", ""))
    order := datasource.NewOrder("begin_ts", "asc")

    rs := datasource.ListRequest(p.db, conditions, order, nil)
    index := make([]int, len(rs))

    for i, r := range rs {
        if r.ParentUuid == "" {
            index[i] = i
            continue
        }

        pi := datasource.FindRequestParent(rs, r)
        if pi < 0 {
            index[i] = -1
            continue
        }

        index[i] = index[pi]
    }

    // process new request group
    for j, k := range index {
        if j == k && time.Now().UTC().Unix() - rs[j].CreateTs > ProcessWaiting {
            group := getRequestGroup(rs, index, j)
            guuid := rs[j].Uuid
            log.Printf("INFO: Start to process new request group %s\n", guuid)
            p.setRequestGroup(group, guuid)
            p.processRequestGroup(guuid)
        }
    }

    // update legacy request group
    updateList := make([]string, 0)
    for j, k := range index {
        if k != -1 {
            continue
        }

        req := datasource.GetRequest(p.db, rs[j].ParentUuid)
        if req != nil && req.GroupUuid != "" {
            reqs := make([]*datasource.Request, 0)
            reqs = append(reqs, rs[j])
            log.Printf("INFO: Add request %s to group %s\n", rs[j].Uuid, req.GroupUuid)
            p.setRequestGroup(reqs, req.GroupUuid)
            insertSet(&updateList, req.GroupUuid)
            continue
        }

        if time.Now().UTC().Unix() - rs[j].CreateTs > ProcessTimeout {
            log.Printf("INFO: Delete request %s which is out of date\n", rs[j].Uuid)
            rs[j].Delete(p.db)
            continue
        }
    }

    for _, guuid := range updateList {
        log.Printf("INFO: Start to process legacy request group %s\n", guuid)
        rg := datasource.GetRequestGroup(p.db, guuid)
        if rg != nil {
            log.Printf("INFO: Delete legacy request group %s\n", rg.Uuid)
            rg.Delete(p.db)
            conditions := make([]*datasource.Condition, 0)
            conditions = append(conditions, datasource.NewCondition("invoke_chain_id", "=", rg.InvokeChainId))
            if len(datasource.ListRequestGroup(p.db, conditions, nil, nil)) == 0 {
                log.Printf("INFO: Delete invoke chain %d since no such request group\n", rg.InvokeChainId)
                datasource.GetInvokeChain(p.db, rg.InvokeChainId).Delete(p.db)
            }
        }

        p.processRequestGroup(guuid)
    }
}

func (p *Project) setRequestGroup(group []*datasource.Request, guuid string) {
    for _, r := range group {
        r.GroupUuid = guuid
        r.Update(p.db)
    }
}

func (p *Project) processRequestGroup(guuid string) {
    conditions := make([]*datasource.Condition, 0)
    conditions = append(conditions, datasource.NewCondition("group_uuid", "=", guuid))
    order := datasource.NewOrder("begin_ts", "asc")
    group := datasource.ListRequest(p.db, conditions, order, nil)
    root := datasource.GetRequest(p.db, guuid)

    log.Printf("Processing request group, uuid = %s\n", guuid)

    seq := make([]*datasource.Request, 0)
    parentSeq := make([]int, 0)
    seq = append(seq, root)
    parentSeq = append(parentSeq, -1)
    for i := 0; i < len(seq); i++ {
        children := datasource.FindRequestChildren(group, seq[i])
        for _, v := range children {
            seq = append(seq, v)
            parentSeq = append(parentSeq, i)
        }
    }

    uuidSeq := make([]string, 0)
    for _, r := range seq {
        uuidSeq = append(uuidSeq, r.Uuid)
    }

    typeSeq := make([]*datasource.RequestType, 0)
    for _, r := range seq {
        rt := new(datasource.RequestType)
        rt.Service = r.Service
        rt.Category = r.Category
        typeSeq = append(typeSeq, rt)
    }

    log.Printf("Uuid sequence, %v\n", uuidSeq)
    log.Printf("Request type sequence, %v\n", typeSeq)
    log.Printf("Parent index sequence %v\n", parentSeq)

    rg := new(datasource.RequestGroup)
    rg.Uuid = root.Uuid
    rg.RequestUuids = uuidSeq
    rg.ParentsIndex = parentSeq

    ic := datasource.GetInvokeChainByValues(p.db, typeSeq, parentSeq)
    if ic != nil {
        rg.InvokeChainId = ic.Id
    } else {
        c := new(datasource.InvokeChain)
        c.Header = datasource.RequestType{
            Service: root.Service,
            Category: root.Category,
        }
        c.RequestTypes = typeSeq
        c.ParentsIndex = parentSeq
        c.Save(p.db)
        rg.InvokeChainId = c.Id
    }

    rg.Save(p.db)
}

////////////////////////////////////////////////////////////////////////////////

func getRequestGroup(rs []*datasource.Request, index []int, value int) []*datasource.Request {
    group := make([]*datasource.Request, 0)
    for i, v := range index {
        if v == value {
            group = append(group, rs[i])
        }
    }

    return group
}

func insertSet(set *[]string, ele string) {
    for _, s := range *set {
        if s == ele {
            return
        }
    }

    *set = append(*set, ele)
}