package webface;

import (
  "fmt"
  "time"
  "github.com/satori/go.uuid"
)

type SessionMeta struct {
  clientId string;
  username string;
  ip       string;
  created  string;
}

type SessionMgr map[string]interface{};
type SessionObj map[string]string;

// TODO: to persist session data
// import "github.com/srinathgs/mysqlstore";
var __webfaceSession SessionMgr = make(SessionMgr);

func __webfaceSessionSanityCheck (m *SessionMeta) (*SessionObj, int) {
  obj := __webfaceSession[m.clientId];
  if (obj == nil) {
    return nil, -1;
  }
  s := obj.(*SessionObj);
  if ((*s)["__webface_cid"] != m.clientId) {
    return nil, 1;
  }
  if ((*s)["__webface_username"] != m.username) {
    return nil, 2;
  }
  if ((*s)["__webface_ip"] != m.ip) {
    return nil, 3;
  }
  if ((*s)["__webface_created"] != m.created) {
    return nil, 4;
  }
  return s, 0;
}

func SessionBegin (username string, ip string) *SessionMeta {
  meta := new(SessionMeta);
  meta.clientId = uuid.NewV4().String();
  meta.username = username;
  meta.ip = ip;
  meta.created = time.Now().String();
  session := new(SessionObj);
  (*session) = make(SessionObj);
  (*session)["__webface_cid"] = meta.clientId;
  (*session)["__webface_username"] = meta.username;
  (*session)["__webface_ip"] = meta.ip;
  (*session)["__webface_created"] = meta.created;
  __webfaceSession[meta.clientId] = session;
  return meta;
}

func SessionEnd (meta *SessionMeta) int {
  var  errcode int;
  _, errcode = __webfaceSessionSanityCheck(meta);
  if (errcode != 0) {
    return errcode;
  }
  delete(__webfaceSession, meta.clientId);
  return 0;
}

func Session (meta *SessionMeta) (*SessionObj, int) {
  var session *SessionObj;
  var  errcode int;
  session, errcode = __webfaceSessionSanityCheck(meta);
  if (errcode != 0) {
    return nil, errcode;
  }
  return session, 0; // exists
}

func SessionGet (s *SessionObj, key string) string {
  return (*s)[key];
}

func SessionSet (s *SessionObj, key string, val string) string {
  old := (*s)[key];
  (*s)[key] = val;
  return old;
}

func TestMain () {
  m := SessionBegin("fakeuser", "127.0.0.1");
  s, _ := Session(m);
  fmt.Println(SessionGet(s, "__webface_username"));
  fmt.Println(SessionSet(s, "__webface_ip", "192.168.1.1"));
  fmt.Println("end session: ", SessionEnd(m));
  fmt.Println(SessionSet(s, "__webface_ip", "127.0.0.1"));
  fmt.Println("end session: ", SessionEnd(m));
}
