package model

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
