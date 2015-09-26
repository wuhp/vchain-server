package model

import (
    "fmt"
    "strings"
)

////////////////////////////////////////////////////////////////////////////////

type Order struct {
    Columns string
    Sequence string
}

func NewOrder(cols, seq string) *Order {
    return &Order{
        Columns: cols,
        Sequence: seq,
    }
}

////////////////////////////////////////////////////////////////////////////////

type Paging struct {
    Offset int
    Size   int
}

func NewPaging(offset, size int) *Paging {
    return &Paging{
        Offset: offset,
        Size: size,
    }
}

////////////////////////////////////////////////////////////////////////////////

type Condition struct {
    Key   string
    Op    string
    Value interface{}
}

func NewCondition(key, op string, value interface{}) *Condition {
    return &Condition{
        Key: key,
        Op: op,
        Value: value,
    }
}

////////////////////////////////////////////////////////////////////////////////

func GenerateOrderSql(o *Order) string {
    if o != nil {
        return fmt.Sprintf(" ORDER BY %s %s", o.Columns, o.Sequence)
    }

    return ""
}

func GenerateLimitSql(p *Paging) string {
    if p != nil {
        return fmt.Sprintf(" LIMIT %d, %d", p.Offset, p.Size)
    }

    return ""
}

func GenerateWhereSql(cs []*Condition) (string, []interface{}) {
    if len(cs) > 0 {
        ks := make([]string, 0)
        vs := make([]interface{}, 0)
        for _, c := range cs {
            ks = append(ks, fmt.Sprintf("%s%s?", c.Key, c.Op))
            vs = append(vs, c.Value)
        }
        return " WHERE " + strings.Join(ks, " and "), vs
    }

    return "", nil
}
