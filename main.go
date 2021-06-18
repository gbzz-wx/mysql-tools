package main

import (
  "database/sql"
  "encoding/json"
  "flag"
  "fmt"
  _ "github.com/go-sql-driver/mysql"
)

func getDb(link string) *sql.DB {
  db, err := sql.Open("mysql", link)
  if err != nil {
    fmt.Printf("连接mysql异常:%e", err)
  }
  return db
}
type RowData map[string]interface{}
func main() {
  var sql string
  var link string
  flag.StringVar(&link, "link", "", `用户名:密码@tcp(地址:端口)/数据库名`)
  flag.StringVar(&sql, "sql", "", "sql")
  flag.Parse()

  db := getDb(link)
  defer db.Close()

  result,err := QueryAll(db, sql)
  if err != nil{
    fmt.Print("连接mysql异常:", err)
  }
  for i := range result {
    b, _ := json.Marshal(result[i])
    fmt.Println(string(b))
  }

}

//查询所有
func QueryAll(db *sql.DB, sqlStr string) ([]*RowData, error) {
  rows, err := db.Query(sqlStr)
  if err != nil {
    return nil, err
  }
  //函数结束释放链接
  defer rows.Close()
  //读出查询出的列字段名
  cols, _ := rows.Columns()
  //values是每个列的值，这里获取到byte里
  values := make([][]byte, len(cols))
  //query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
  scans := make([]interface{}, len(cols))
  //让每一行数据都填充到[][]byte里面,狸猫换太子
  for i := range values {
    scans[i] = &values[i]
  }
  results := make([]*RowData, 0, 10)
  for rows.Next() {
    err := rows.Scan(scans...)
    if err != nil {
      return nil, err
    }
    row := make(RowData, 10)
    for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
      key := cols[k]
      row[key] = string(v)
    }
    results = append(results, &row)
  }
  //返回数据
  return results, nil
}
