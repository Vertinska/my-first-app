package main

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/Vertinska/my-first-app/internal/repository"
)

type sqliteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) repository.ProductRepository {
    return &sqliteRepository{db: db}
}

func (r *sqliteRepository) GetProduct(ctx context.Context, id int) (*repository.Product, error) {
    row := r.db.QueryRowContext(ctx, "SELECT id, model, company, price FROM products WHERE id = ?", id)
    p := &repository.Product{}
    err := row.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("product with id %d not found", id)
    }
    return p, err
}

func (r *sqliteRepository) ListProducts(ctx context.Context) ([]repository.Product, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT id, model, company, price FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []repository.Product
    for rows.Next() {
        var p repository.Product
        if err := rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, rows.Err()
}

func (r *sqliteRepository) SaveProduct(ctx context.Context, product repository.Product) error {
    result, err := r.db.ExecContext(ctx,
        "INSERT INTO products (model, company, price) VALUES (?, ?, ?)",
        product.Model, product.Company, product.Price)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    // Обновляем ID в переданной структуре
    product.ID = int(id)
    return nil
}

func (r *sqliteRepository) GetExpensiveProducts(ctx context.Context, minPrice int) ([]repository.Product, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT id, model, company, price FROM products WHERE price >= ?", minPrice)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []repository.Product
    for rows.Next() {
        var p repository.Product
        if err := rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, rows.Err()
}

func (r *sqliteRepository) DeleteProduct(ctx context.Context, id int) error {
    result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product with id %d not found", id)
    }
    
    return nil
}
