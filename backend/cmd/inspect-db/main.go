// inspect-db prints schema metadata only; it never writes database records.
package main

import (
 "bufio"
 "database/sql"
 "fmt"
 "log"
 "net"
 "net/url"
 "os"
 "strings"
 "time"
 mysql "github.com/go-sql-driver/mysql"
 _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
 f, err := os.Open(".env"); if err != nil { log.Fatal(err) }; defer f.Close()
 s := bufio.NewScanner(f); for s.Scan() { k,v,ok := strings.Cut(s.Text(), "="); if ok && !strings.HasPrefix(strings.TrimSpace(k), "#") { os.Setenv(strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), "\"'")) } }; if err := s.Err(); err != nil { log.Fatal(err) }
 for _, prefix := range []string{"DB_MYSQL", "DB_PG", "DB_PG2"} {
  var driver, dsn string
  if prefix == "DB_MYSQL" { c := mysql.NewConfig(); c.User=os.Getenv(prefix+"_USER"); c.Passwd=os.Getenv(prefix+"_PASS"); c.Net="tcp"; c.Addr=net.JoinHostPort(os.Getenv(prefix+"_HOST"),os.Getenv(prefix+"_PORT")); c.DBName=os.Getenv(prefix+"_NAME"); c.Timeout=5*time.Second; driver="mysql"; dsn=c.FormatDSN() } else { u:=url.URL{Scheme:"postgres",User:url.UserPassword(os.Getenv(prefix+"_USER"),os.Getenv(prefix+"_PASS")),Host:net.JoinHostPort(os.Getenv(prefix+"_HOST"),os.Getenv(prefix+"_PORT")),Path:os.Getenv(prefix+"_NAME")}; q:=u.Query(); q.Set("sslmode","disable"); q.Set("connect_timeout","5"); u.RawQuery=q.Encode(); driver="pgx"; dsn=u.String() }
  db,err:=sql.Open(driver,dsn); if err!=nil { log.Fatal(err) }; defer db.Close()
  fmt.Println(prefix)
  if err=db.Ping(); err!=nil { fmt.Println("connection failed:",err); continue }
  query:=`SELECT table_name, column_name, column_type, is_nullable, COALESCE(column_default,'<NULL>'), extra FROM information_schema.columns WHERE table_schema=DATABASE() ORDER BY table_name,ordinal_position`
  if prefix=="DB_PG" { query=`SELECT table_schema,table_name,column_name,data_type FROM information_schema.columns WHERE table_schema='master' AND table_name IN ('branch','employee') ORDER BY table_name,ordinal_position` }
  if prefix=="DB_PG2" { query=`SELECT table_schema,table_name,column_name,data_type FROM information_schema.columns WHERE table_schema NOT IN ('pg_catalog','information_schema') ORDER BY table_schema,table_name,ordinal_position` }
  rows,err:=db.Query(query); if err!=nil {log.Fatal(err)}; cols,_:=rows.Columns(); for rows.Next(){ vals:=make([]sql.NullString,len(cols)); ptrs:=make([]any,len(cols)); for i:=range vals {ptrs[i]=&vals[i]}; if err:=rows.Scan(ptrs...); err!=nil {log.Fatal(err)}; out:=make([]string,len(cols)); for i,v:=range vals {out[i]=v.String}; fmt.Println(strings.Join(out," | ")) }; if err:=rows.Err();err!=nil {log.Fatal(err)}; rows.Close()
 }
}
