package model

import "time"

type Repo struct {
    Id       int64  `json:"id"`
    UserId   int64  `json:"user_id"`
    Name     string `json:"name"`
    Hash     string `json:"hash"`
    CreateTs int64  `json:"create_ts"`
}


var TableDefinition_Repo []string = []string {
  "id", "user_id", "name", "hash", "create_ts",
}

func Map2Repo (obj map[string]interface{}) *Repo {
  res := new(Repo)
  res.Id = obj["id"].(int64)
  res.UserId = obj["user_id"].(int64)
  res.Name = obj["name"].(string)
  res.Hash = obj["hash"].(string)
  res.CreateTs = obj["create_ts"].(int64)
  return res
}

func Repo2Map (obj *Repo) map[string]interface{} {
  res := make(map[string]interface{})
  res["id"] = obj.Id
  res["user_id"] = obj.UserId
  res["name"] = obj.Name
  res["hash"] = obj.Hash
  res["create_ts"] = obj.CreateTs
  return res
}

func RepoMakeValues (repo *Repo) []interface{} {
  values := make([]interface{}, 4)
  values[0] = repo.UserId
  values[1] = repo.Name
  values[2] = repo.Hash
  values[3] = repo.CreateTs
  return values
}

func RepoGetOne(id int64) map[string]interface{} {
  sql := "SELECT " + TableColumns(TableDefinition_Repo) + " FROM repos WHERE id=?"
  return DBGetOne(db, sql, id, TableDefinition_Repo)
}

func RepoGetList(conditions string, args []interface{}) []map[string]interface{} {
  sql := "SELECT " + TableColumns(TableDefinition_Repo) + " FROM repos " + conditions
  return DBGetList(db, sql, args, TableDefinition_Repo)
}

func RepoCreate(repo *Repo) *Repo {
  sql := "INSERT INTO repos (user_id, name, hash, create_ts) VALUES (?, ?, ?, ?)"
  repo.CreateTs = time.Now().UTC().Unix()
  id := DBInsert(db, sql, RepoMakeValues(repo))
  repo.Id = id
  return repo
}

func RepoUpdate(repo *Repo) *Repo {
  sql := "UPDATE repos SET user_id=?, name=?, hash=?, create_ts=? WHERE id=?"
  DBUpdate(db, sql, repo.Id, RepoMakeValues(repo))
  return repo
}

func RepoRemove(id int64) {
  sql := "DELETE FROM repos WHERE id=?"
  DBRemove(db, sql, id)
}

func (repo *Repo) ResetHash() {
  repo.Hash = generateHash()
  RepoUpdate(repo)
}
