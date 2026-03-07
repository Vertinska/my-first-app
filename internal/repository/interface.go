package repository

import "context"

//go:generate mockgen -source=./interface.go -destination=../../mocks/repository_mock.go -package=mocks
type ProductRepository interface {
    GetProduct(ctx context.Context, id int) (*Product, error)
    ListProducts(ctx context.Context) ([]Product, error)
    SaveProduct(ctx context.Context, product Product) error
    GetExpensiveProducts(ctx context.Context, minPrice int) ([]Product, error)
}

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}
