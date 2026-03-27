package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    "github.com/Vertinska/my-first-app/internal/repository"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    ctx := context.Background()
    fmt.Println("🚀 Начинаем работу с базой данных...")

    // 1. Подключение к базе
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal("❌ Ошибка подключения:", err)
    }
    defer db.Close()
    fmt.Println("✅ Подключение к базе данных установлено")

    // 2. Создаём репозиторий
    repo := NewSQLiteRepository(db)

    // 3. Создание таблицы
    createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`

    _, err = db.ExecContext(ctx, createTableSQL)
    if err != nil {
        log.Fatal("❌ Ошибка создания таблицы:", err)
    }
    fmt.Println("✅ Таблица products проверена/создана")

    // 4. Добавление тестовых данных
    testProducts := []repository.Product{
        {Model: "iPhone X", Company: "Apple", Price: 72000},
        {Model: "Galaxy S20", Company: "Samsung", Price: 65000},
        {Model: "Xiaomi Mi 11", Company: "Xiaomi", Price: 40000},
    }

    fmt.Println("\n📦 Добавление тестовых товаров:")
    for _, p := range testProducts {
        err := repo.SaveProduct(ctx, p)
        if err != nil {
            log.Printf("❌ Ошибка вставки %s: %v\n", p.Model, err)
        } else {
            fmt.Printf("  ✅ Добавлен: %s (%s) - %d руб.\n", p.Model, p.Company, p.Price)
        }
    }

    // 5. Получение всех данных
    fmt.Println("\n📦 Список всех товаров:")
    allProducts, err := repo.ListProducts(ctx)
    if err != nil {
        log.Fatal("❌ Ошибка получения списка:", err)
    }
    if len(allProducts) == 0 {
        fmt.Println("  В базе нет товаров")
    } else {
        for _, p := range allProducts {
            fmt.Printf("  %d. %s (%s) - %d руб.\n", p.ID, p.Model, p.Company, p.Price)
        }
    }

    // 6. Проверка поиска по цене
    fmt.Println("\n💰 Товары дороже 50000 руб.:")
    expensive, err := repo.GetExpensiveProducts(ctx, 50000)
    if err != nil {
        log.Fatal("❌ Ошибка запроса:", err)
    }
    if len(expensive) == 0 {
        fmt.Println("  Нет товаров дороже 50000 руб.")
    } else {
        for _, p := range expensive {
            fmt.Printf("  %d. %s - %d руб.\n", p.ID, p.Model, p.Price)
        }
    }

    // 7. Получение одной записи
    fmt.Println("\n🔍 Товар с ID=1:")
    oneProduct, err := repo.GetProduct(ctx, 1)
    if err != nil {
        fmt.Printf("  ❌ Товар с ID=1 не найден: %v\n", err)
    } else {
        fmt.Printf("  ✅ %d. %s (%s) - %d руб.\n", oneProduct.ID, oneProduct.Model, oneProduct.Company, oneProduct.Price)
    }

    fmt.Println("\n✨ Проверка завершена!")
}
