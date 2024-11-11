package main

import (
  "database/sql"
  "fmt"
  "log"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  "os"
)

func main() {
  dbUser := os.Getenv("MYSQL_USER")
  dbName := os.Getenv("MYSQL_DATABASE")
  dbHost := "mysql"
  env := os.Getenv("ENV")

  dbPassword := os.Getenv("MYSQL_PASSWORD")
  if dbPassword == "" {
      log.Fatal("MySQL password not set")
  }

  dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT NOW()")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    defer rows.Close()

    var currentTime string
    if rows.Next() {
      if err := rows.Scan(&currentTime); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
    }
    fmt.Fprintf(w, "Current time: %s. Env: %s", currentTime, env)
  })

  log.Println("Server started on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
