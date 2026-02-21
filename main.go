package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    fmt.Println("🚀 Начинаем работу с базой данных...")
    
    // 1. Подключение к базе
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal("❌ Ошибка подключения:", err)
    }
    defer db.Close()
    fmt.Println("✅ Подключение к базе данных установлено")

    // 2. Создание таблицы
    createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`
    
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal("❌ Ошибка создания таблицы:", err)
    }
    fmt.Println("✅ Таблица products проверена/создана")

    // 3. Добавление тестовых данных
    products := []struct {
        model   string
        company string
        price   int
    }{
        {"iPhone X", "Apple", 72000},
        {"Galaxy S20", "Samsung", 65000},
        {"Xiaomi Mi 11", "Xiaomi", 40000},
    }

    for _, p := range products {
        result, err := db.Exec(
            "INSERT INTO products (model, company, price) VALUES (?, ?, ?)",
            p.model, p.company, p.price,
        )
        if err != nil {
            log.Printf("❌ Ошибка вставки %s: %v\n", p.model, err)
        } else {
            id, _ := result.LastInsertId()
            fmt.Printf("✅ Добавлен товар: %s (ID: %d)\n", p.model, id)
        }
    }

    // 4. Получение всех данных
    fmt.Println("\n📦 Список всех товаров:")
    rows, err := db.Query("SELECT * FROM products")
    if err != nil {
        log.Fatal("❌ Ошибка запроса:", err)
    }
    defer rows.Close()

    for rows.Next() {
        p := Product{}
        err := rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
        if err != nil {
            log.Println("❌ Ошибка сканирования:", err)
            continue
        }
        fmt.Printf("  %d. %s (%s) - %d руб.\n", p.ID, p.Model, p.Company, p.Price)
    }

    // 5. Проверка поиска по цене
    fmt.Println("\n💰 Товары дороже 50000 руб.:")
    expensiveRows, err := db.Query("SELECT * FROM products WHERE price > ?", 50000)
    if err != nil {
        log.Fatal("❌ Ошибка запроса:", err)
    }
    defer expensiveRows.Close()

    for expensiveRows.Next() {
        p := Product{}
        err := expensiveRows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
        if err != nil {
            continue
        }
        fmt.Printf("  %d. %s - %d руб.\n", p.ID, p.Model, p.Price)
    }

    // 6. Получение одной записи
    fmt.Println("\n🔍 Товар с ID=1:")
    row := db.QueryRow("SELECT * FROM products WHERE id = ?", 1)
    p := Product{}
    err = row.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("  Товар с ID=1 не найден")
        } else {
            log.Println("❌ Ошибка:", err)
        }
    } else {
        fmt.Printf("  %d. %s (%s) - %d руб.\n", p.ID, p.Model, p.Company, p.Price)
    }

    fmt.Println("\n✨ Проверка завершена!")
}
