package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/Vertinska/my-first-app/internal/repository"
)

type ItemHandler struct {
    repo repository.ProductRepository
}

func NewItemHandler(repo repository.ProductRepository) *ItemHandler {
    return &ItemHandler{repo: repo}
}

// CreateItem — POST /items
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
    var product repository.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.repo.SaveProduct(r.Context(), product); err != nil {
        http.Error(w, "Failed to create item", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

// GetAllItems — GET /items
func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
    products, err := h.repo.ListProducts(r.Context())
    if err != nil {
        http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

// GetItemByID — GET /items/{id}
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    idStr := pathParts[2]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid item ID", http.StatusBadRequest)
        return
    }

    product, err := h.repo.GetProduct(r.Context(), id)
    if err != nil {
        http.Error(w, "Item not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}

// DeleteItem — DELETE /items/{id}
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    idStr := pathParts[2]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid item ID", http.StatusBadRequest)
        return
    }

    if err := h.repo.DeleteProduct(r.Context(), id); err != nil {
        http.Error(w, "Failed to delete item", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
