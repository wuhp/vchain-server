package parser

import (
    "log"
    "time"

    "vchaind/model"
)

func MainLoop() {
    for {
        conditions := make([]*model.Condition, 0)
        conditions = append(conditions, model.NewCondition("group_uuid", "=", ""))
        conditions = append(conditions, model.NewCondition("end_ts", "!=", "0"))
        order := model.NewOrder("begin_ts", "asc")
        process(model.ListRequest(conditions, order, nil))
        time.Sleep(5 * time.Second)
    }
}

func getRequestGroup(rs []*model.Request, index []int, value int) []*model.Request {
    group := make([]*model.Request, 0)
    for i, v := range index {
        if v == value {
            group = append(group, rs[i])
        }
    }

    return group
}

func processRequestGroup(group []*model.Request) {
    root := group[0]
    log.Printf("Processing request group, uuid = %s\n", root.Uuid)

    seq := make([]*model.Request, 0)
    parentSeq := make([]int, 0)
    seq = append(seq, root)
    parentSeq = append(parentSeq, -1)
    for i := 0; i < len(seq); i++ {
        children := model.FindRequestChildren(group, seq[i])
        for _, v := range children {
            seq = append(seq, v)
            parentSeq = append(parentSeq, i)
        }
    }

    uuidSeq := make([]string, 0)
    for _, r := range seq {
        uuidSeq = append(uuidSeq, r.Uuid)
    }

    typeSeq := make([]*model.RequestType, 0)
    for _, r := range seq {
        rt := new(model.RequestType)
        rt.Service = r.Service
        rt.Category = r.Category
        typeSeq = append(typeSeq, rt)
    }

    log.Printf("Uuid sequence, %v\n", uuidSeq)
    log.Printf("Request type sequence, %v\n", typeSeq)
    log.Printf("Parent index sequence %v\n", parentSeq)

    rg := new(model.RequestGroup)
    rg.Uuid = root.Uuid
    rg.RequestUuids = uuidSeq
    rg.ParentsIndex = parentSeq

    ic := model.GetInvokeChainByValues(typeSeq, parentSeq)
    if ic != nil {
        rg.InvokeChainId = ic.Id
    } else {
        c := new(model.InvokeChain)
        c.Header = model.RequestType{
            Service: root.Service,
            Category: root.Category,
        }
        c.RequestTypes = typeSeq
        c.ParentsIndex = parentSeq
        c.Save()
        rg.InvokeChainId = c.Id
    }

    rg.Save()

    for _, r := range group {
        r.GroupUuid = root.Uuid
        r.Update()
    }
}

func process(rs []*model.Request) {
    index := make([]int, len(rs))
    for i, r := range rs {
        if r.ParentUuid == "" {
            index[i] = i
            continue
        }

        pi := model.FindRequestParent(rs, r)
        if pi < 0 {
            index[i] = -1
            continue
        }

        index[i] = index[pi]
    }

    for j, k := range index {
        if j == k {
            group := getRequestGroup(rs, index, j)
            processRequestGroup(group)
        }
    }
}
