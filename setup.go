package main

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal("Ошибка создания таблицы:", err)
    }
    log.Println("Таблица успешно создана")
}
