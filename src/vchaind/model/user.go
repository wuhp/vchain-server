package model

import "time"

type User struct {
    Id            int64  `json:"id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    Password      string `json:"-"`
    CreateTs      int64  `json:"-"`
    LastLoginTs   int64  `json:"-"`
}

var TableDefinition_User []string = []string {
  "id", "name", "email", "email_verified", "password", "create_ts", "last_login_ts",
}

func internalHash(plain string) string {
  // TODO: HASH(plain)
  return plain
}

func Map2User (obj map[string]interface{}) *User {
  res := new(User)
  res.Id = obj["id"].(int64)
  res.Name = obj["name"].(string)
  res.Email = obj["email"].(string)
  res.EmailVerified = obj["email_verified"].(bool)
  res.Password = obj["password"].(string)
  res.CreateTs = obj["create_ts"].(int64)
  res.LastLoginTs = obj["last_login_ts"].(int64)
  return res
}

func User2Map (obj *User) map[string]interface{} {
  res := make(map[string]interface{})
  res["id"] = obj.Id
  res["name"] = obj.Name
  res["email"] = obj.Email
  res["email_verified"] = obj.EmailVerified
  res["password"] = obj.Password
  res["create_ts"] = obj.CreateTs
  res["last_login_ts"] = obj.LastLoginTs
  return res
}

func UserMakeValues (user *User) []interface{} {
  values := make([]interface{}, 6)
  values[0] = user.Name
  values[1] = user.Email
  values[2] = user.EmailVerified
  values[3] = user.Password
  values[3] = user.CreateTs
  values[3] = user.LastLoginTs
  return values
}

func UserGetOne(id int64) map[string]interface{} {
  sql := "SELECT " + TableColumns(TableDefinition_User) + " FROM users WHERE id=?"
  return DBGetOne(db, sql, id, TableDefinition_User)
}

func UserGetList(conditions string, args []interface{}) []map[string]interface{} {
  sql := "SELECT " + TableColumns(TableDefinition_User) + " FROM users " + conditions
  arr := DBGetList(db, sql, args, TableDefinition_User)
  for _, user := range arr {
    user["password"] = ""
  }
  return arr
}

func UserCreate(user *User) *User {
  sql := `INSERT INTO
          users  (name, email, email_verified, password, create_ts, last_login_ts)
          VALUES (?, ?, ?, ?, ?, ?)`
  user.CreateTs = time.Now().UTC().Unix()
  user.Password = internalHash(user.Password)
  id := DBInsert(db, sql, UserMakeValues(user))
  user.Id = id
  return user
}

func UserUpdate(user *User) *User {
  sql := `UPDATE users
          SET name=?, email=?, email_verified=?,
              password=?, create_ts=? last_login_ts =?
          WHERE id=?`
  if user.Password != "" {
    // if password field is not empty, change password
    user.Password = internalHash(user.Password)
  }
  DBUpdate(db, sql, user.Id, UserMakeValues(user))
  return user
}

func UserRemove(id int64) {
  sql := "DELETE FROM users WHERE id=?"
  DBRemove(db, sql, id)
}

func (user *User) Login(username, password string) *User {
  user.Id = -1
  password = internalHash(password)
  values := make([]interface{}, 2)
  values[0] = username
  values[1] = password
  arr := UserGetList("WHERE name=? AND password=?", values)
  if len(arr) == 0 {
    return nil;
  }
  if len(arr) > 1 {
    // should not be
    return nil;
  }
  obj := arr[0]
  user.Id = obj["id"].(int64)
  user.Name = obj["name"].(string)
  user.Email = obj["email"].(string)
  user.EmailVerified = obj["email_verified"].(bool)
  user.Password = ""
  user.CreateTs = obj["create_ts"].(int64)
  user.LastLoginTs = time.Now().UTC().Unix()
  UserUpdate(user)
  return user
}
