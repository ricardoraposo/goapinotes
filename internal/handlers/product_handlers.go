package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ricardoraposo/api-again/internal/database"
	"github.com/ricardoraposo/api-again/internal/dto"
	"github.com/ricardoraposo/api-again/internal/entity"
	entityPkg "github.com/ricardoraposo/api-again/pkg/entity"
)


type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p dto.CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := entity.NewProduct(p.Name, p.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page, err := strconv.Atoi(pageStr)
    if err != nil {
        page = 0
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        limit = 0
    }
    sort := r.URL.Query().Get("sort")

    products, err := h.ProductDB.FindAll(page, limit, sort)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    product, err := h.ProductDB.FindById(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var product entity.Product
    err := json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    product.ID, err = entityPkg.ParseID(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    _, err = h.ProductDB.FindById(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    err = h.ProductDB.Update(&product)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler ) DeleteProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    _, err := h.ProductDB.FindById(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    err = h.ProductDB.Delete(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
