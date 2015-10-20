package model

import (
    "fmt"
    "strings"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDatabase(host, port, user, password, database string) (err error) {
  uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

  if db, err = sql.Open("mysql", uri); err != nil {
      return
  }

  err = db.Ping()
  return
}

func checkError(err error, function string) {
  if err != nil {
    panic(err)
  }
}

func TableColumns(columns []string) string {
  return strings.Join(columns[:], ",")
}

func DBExec(db *sql.DB, sql string, values []interface{}) sql.Result {
  stmt, err := db.Prepare(sql)
  checkError(err, "DbExec:[s]=>Prepare")
  defer stmt.Close()
  res, err := stmt.Exec(values...)
  checkError(err, "DbExec:Prepare=>Exec")
  return res
}

func DBInsert(db *sql.DB, sql string, values []interface{}) int64 {
  id, err := DBExec(db, sql, values).LastInsertId()
  checkError(err, "DBInsert:DBExec=>LastInsertId")
  return id
}

func DBUpdate(db *sql.DB, sql string, id int64, values []interface{}) {
  values = append(values, id)
  DBExec(db, sql, values)
}

func DBRemove(db *sql.DB, sql string, id int64) {
  values := make([]interface{}, 1, 1)
  values[0] = id
  DBExec(db, sql, values)
}

func DBRowRead(row *sql.Row, columns []string) map[string]interface{} {
  dict := make(map[string]interface{})
  val := make([]interface{}, len(columns))
  ptr := make([]interface{}, len(columns))
  for i, _ := range columns {
    ptr[i] = &val[i]
  }
  err := row.Scan(ptr...)
  checkError(err, "DBRowRead")
  for i, column := range columns {
    dict[column] = val[i]
  }
  return dict
}

func DBRowsLineRead(rows *sql.Rows, columns []string) map[string]interface{} {
  dict := make(map[string]interface{})
  val := make([]interface{}, len(columns))
  ptr := make([]interface{}, len(columns))
  for i, _ := range columns {
    ptr[i] = &val[i]
  }
  err := rows.Scan(ptr...)
  checkError(err, "DBRowsLineRead")
  for i, column := range columns {
    dict[column] = val[i]
  }
  return dict
}

func DBGetOne(db *sql.DB, sql string, id int64, columns []string) map[string]interface{} {
  values := make([]interface{}, 1, 1)
  values[0] = id
  res := db.QueryRow(sql, values...)
  return DBRowRead(res, columns)
}

func DBGetList(db *sql.DB, sql string, values []interface{}, columns []string) []map[string]interface{} {
  res, err := db.Query(sql, values...)
  checkError(err, "DBGetList")
  defer res.Close()
  arr := make([]map[string]interface{}, 0)
  for res.Next() {
    arr = append(arr, DBRowsLineRead(res, columns))
  }
  return arr
}
