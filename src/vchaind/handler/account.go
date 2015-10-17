package handler

import (
  "net/http"
  "encoding/json"
  "time"
  "sync"

  "github.com/satori/go.uuid"

  "vchaind/model"
)

var mutex = &sync.Mutex{}
var shield int64 = 0
var ram_userdb map[int64]*model.User = make(map[int64]*model.User)

func lookup (name, password string) *model.User {
  for _, val := range ram_userdb {
    if val.Name == name {
      if password == "" {
        return val
      }
      if val.Password == password {
        return val
      }
    }
  }
  return nil
}

func Register(w http.ResponseWriter, r *http.Request) {
  mutex.Lock()
  defer mutex.Unlock()
  // secure constraint: 1 sec per register
  cur := time.Now().UnixNano() / 1000000000
  if (cur < shield) {
    http.Error(w, "VCHAIN Server Shield Protect: try again later", 400)
    return
  }

  err := r.ParseForm()
  if err != nil {
    http.Error(w, "Cannot parse POST form", 400)
    return
  }

  name := r.PostForm.Get("name")
  email := r.PostForm.Get("email")
  if name == "" || email == "" {
    http.Error(w, "Invalid form data", 400)
    return
  }

  u0 := lookup(name, "")
  if u0 != nil {
    http.Error(w, "User registered", 400)
    return
  }

  u := new(model.User)
  u.Id = cur
  u.Name = name
  u.Email = email
  u.EmailVerified = false
  u.Password = uuid.NewV4().String()
  u.CreateTs = cur
  u.LastLoginTs = 0
  ram_userdb[cur] = u

  shield = cur + 1
  result := make(map[string]interface{})
  result["user"] = u
  result["key"] = u.Password
  json.NewEncoder(w).Encode(result)
}

func Login(w http.ResponseWriter, r *http.Request) {
}
