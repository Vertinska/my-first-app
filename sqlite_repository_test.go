package main

import (
    "context"
    "database/sql"
    "testing"

    "github.com/Vertinska/my-first-app/internal/repository"
    _ "github.com/mattn/go-sqlite3"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)

    createTableSQL := `CREATE TABLE products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`
    _, err = db.Exec(createTableSQL)
    require.NoError(t, err)

    return db
}

func TestSQLiteRepository_SaveAndGetProduct(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewSQLiteRepository(db)
    ctx := context.Background()

    product := repository.Product{
        Model:   "Test Phone",
        Company: "Test Company",
        Price:   50000,
    }

    err := repo.SaveProduct(ctx, product)
    require.NoError(t, err)

    saved, err := repo.GetProduct(ctx, 1)
    require.NoError(t, err)
    require.Equal(t, "Test Phone", saved.Model)
    require.Equal(t, 50000, saved.Price)
}
func TestSQLiteRepository_ListProducts(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewSQLiteRepository(db)
    ctx := context.Background()

    // Добавим тестовые данные
    products := []repository.Product{
        {Model: "Phone 1", Company: "Company A", Price: 10000},
        {Model: "Phone 2", Company: "Company B", Price: 20000},
    }

    for _, p := range products {
        err := repo.SaveProduct(ctx, p)
        require.NoError(t, err)
    }

    // Тестируем ListProducts
    list, err := repo.ListProducts(ctx)
    require.NoError(t, err)
    require.Len(t, list, 2)
}

func TestSQLiteRepository_GetExpensiveProducts(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewSQLiteRepository(db)
    ctx := context.Background()

    // Добавим тестовые данные с разными ценами
    products := []repository.Product{
        {Model: "Cheap", Price: 10000},
        {Model: "Medium", Price: 30000},
        {Model: "Expensive", Price: 60000},
        {Model: "Very Expensive", Price: 90000},
    }

    for _, p := range products {
        err := repo.SaveProduct(ctx, p)
        require.NoError(t, err)
    }

    // Тестируем GetExpensiveProducts с порогом 50000
    expensive, err := repo.GetExpensiveProducts(ctx, 50000)
    require.NoError(t, err)
    require.Len(t, expensive, 2) // Должно быть 2 дорогих товара
    
    for _, p := range expensive {
        require.Greater(t, p.Price, 50000)
    }
}

func TestSQLiteRepository_GetProduct_NotFound(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewSQLiteRepository(db)
    ctx := context.Background()

    // Пытаемся получить несуществующий товар
    product, err := repo.GetProduct(ctx, 999)
    
    // Проверяем, что вернулась ошибка
    require.Error(t, err)
    require.Nil(t, product)
    require.Contains(t, err.Error(), "not found")
}
func TestSQLiteRepository_ListProducts_Error(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Закроем базу данных принудительно, чтобы вызвать ошибку
    db.Close()

    repo := NewSQLiteRepository(db)
    ctx := context.Background()

    // Попытка выполнить запрос к закрытой базе вызовет ошибку
    products, err := repo.ListProducts(ctx)
    
    // Проверяем, что вернулась ошибка
    require.Error(t, err)
    require.Nil(t, products)
}
