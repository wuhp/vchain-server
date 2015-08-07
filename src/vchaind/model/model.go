package model

type Request struct {
    Uuid          string `json:"uuid"`
    ParentUuid    string `json:"parent_uuid"`
    Service       string `json:"service"`
    Category      string `json:"category"`
    SyncOption    string `json:"sync_option"`
    BeginTs       int64  `json:"begin_ts"`
    EndTs         int64  `json:"end_ts"`
    BeginMetadata string `json:"begin_metadata"`
    EndMetadata   string `json:"end_metadata"`
    CreateTs      string `json:"create_ts"`
    UpdateTs      string `json:"update_ts"`
}

type RequestGroup struct {
  Uuid             string
  RequestSeq       string
  RequestParentSeq string
  InvokeChainId    int64
  InProcessOrNot   bool
}

type InvokeChain struct {
  Id               int64
  Header           string
  RequestSeq       string
  RequestParentSeq string
}
