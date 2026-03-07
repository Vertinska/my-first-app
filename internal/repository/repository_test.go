package repository_test

import (
    "context"
    "testing"

    "github.com/Vertinska/my-first-app/internal/repository"
    "github.com/Vertinska/my-first-app/mocks"
    "go.uber.org/mock/gomock"
    "github.com/stretchr/testify/require"
)

func TestGetProduct(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockProductRepository(ctrl)

    expectedProduct := &repository.Product{ID: 1, Model: "iPhone X", Company: "Apple", Price: 72000}
    mockRepo.EXPECT().
        GetProduct(gomock.Any(), 1).
        Return(expectedProduct, nil)

    product, err := mockRepo.GetProduct(context.Background(), 1)
    require.NoError(t, err)
    require.Equal(t, expectedProduct, product)
}

func TestListProducts(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockProductRepository(ctrl)

    expectedProducts := []repository.Product{
        {ID: 1, Model: "iPhone X", Company: "Apple", Price: 72000},
        {ID: 2, Model: "Galaxy S20", Company: "Samsung", Price: 65000},
    }

    mockRepo.EXPECT().
        ListProducts(gomock.Any()).
        Return(expectedProducts, nil)

    products, err := mockRepo.ListProducts(context.Background())
    require.NoError(t, err)
    require.Equal(t, expectedProducts, products)
}

func TestGetExpensiveProducts(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockProductRepository(ctrl)

    expectedProducts := []repository.Product{
        {ID: 1, Model: "iPhone X", Company: "Apple", Price: 72000},
        {ID: 2, Model: "Galaxy S20", Company: "Samsung", Price: 65000},
    }

    mockRepo.EXPECT().
        GetExpensiveProducts(gomock.Any(), 50000).
        Return(expectedProducts, nil)

    products, err := mockRepo.GetExpensiveProducts(context.Background(), 50000)
    require.NoError(t, err)
    require.Equal(t, expectedProducts, products)
}
