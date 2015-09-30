package datasource

import (
    "fmt"
    "strings"
    "strconv"
)

////////////////////////////////////////////////////////////////////////////////

type TimeRange struct {
    Begin int64
    End   int64
}

func NewTimeRange(begin, end int64) *TimeRange {
    return &TimeRange{
        Begin: begin,
        End: end,
    }
}

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

func generateOrderSql(o *Order) string {
    if o != nil {
        return fmt.Sprintf(" ORDER BY %s %s", o.Columns, o.Sequence)
    }

    return ""
}

func generateLimitSql(p *Paging) string {
    if p != nil {
        return fmt.Sprintf(" LIMIT %d, %d", p.Offset, p.Size)
    }

    return ""
}

func generateWhereSql(cs []*Condition) (string, []interface{}) {
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

////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////

func String2strings(s string) []string {
    return strings.Split(s, ",")
}

func Strings2string(ss []string) string {
    return strings.Join(ss, ",")
}

func String2ints(s string) []int {
    ss := strings.Split(s, ",")
    is := make([]int, len(ss))
    for i := 0; i < len(ss); i++ {
        is[i], _ = strconv.Atoi(ss[i])
    }

    return is
}

func Ints2string(is []int) string {
    ss := make([]string, len(is))
    for i := 0; i < len(is); i++ {
        ss[i] = strconv.Itoa(is[i])
    }

    return strings.Join(ss, ",")
}

func String2requestType(s string) *RequestType {
    ss := strings.Split(s, ",")
    
    return &RequestType{
        Service: ss[0],
        Category: ss[1],
    }
}

func RequestType2string(rt *RequestType) string {
    return fmt.Sprintf("%s,%s", rt.Service, rt.Category)
}

func String2requestTypes(s string) []*RequestType {
    ss := strings.Split(s, ";")
    rts := make([]*RequestType, len(ss))
    for i := 0; i < len(ss); i++ {
        rts[i] = String2requestType(ss[i])
    }

    return rts
}

func RequestTypes2string(rts []*RequestType) string {
    ss := make([]string, len(rts))
    for i := 0; i < len(rts); i++ {
        ss[i] = RequestType2string(rts[i])
    }

    return strings.Join(ss, ";")
}

////////////////////////////////////////////////////////////////////////////////

func FindRequestByUuid(rs []*Request, uuid string) int {
    for i, r := range rs {
        if r.Uuid == uuid {
            return i
        }
    }

    return -1
}

func FindRequestParent(rs []*Request, r *Request) int {
    for i, v := range rs {
        if v.Uuid == r.ParentUuid {
            return i
        }
    }

    return -1
}

func FindRequestChildren(rs []*Request, p *Request) []*Request {
    children := make([]*Request, 0)
    for _, r := range rs {
        if r.ParentUuid == p.Uuid {
            children = append(children, r)
        }
    }

    return children
}
