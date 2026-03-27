package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/Vertinska/my-first-app/internal/handler"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Подключение к базе данных
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Создание таблицы если не существует
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT NOT NULL,
        company TEXT NOT NULL,
        price INTEGER NOT NULL
    );`
    if _, err := db.Exec(createTableSQL); err != nil {
        log.Fatal("Failed to create table:", err)
    }

    // Инициализация репозитория
    productRepo := NewSQLiteRepository(db)

    // Создаём обработчик
    itemHandler := handler.NewItemHandler(productRepo)

    // Настройка маршрутов
    http.HandleFunc("POST /items", itemHandler.CreateItem)
    http.HandleFunc("GET /items", itemHandler.GetAllItems)
    http.HandleFunc("GET /items/", itemHandler.GetItemByID)
    http.HandleFunc("DELETE /items/", itemHandler.DeleteItem)

    // Запуск сервера
    log.Println("🚀 Server starting on http://localhost:8080")
    log.Println("📋 Endpoints:")
    log.Println("   POST   /items     - Create item")
    log.Println("   GET    /items     - Get all items")
    log.Println("   GET    /items/{id} - Get item by ID")
    log.Println("   DELETE /items/{id} - Delete item by ID")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
