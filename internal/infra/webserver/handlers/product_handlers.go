package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/reinaldosaraiva/go-api/internal/dto"
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/reinaldosaraiva/go-api/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}
func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// Create Product godoc
// @Summary Create Product
// @Description Create Products
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductDTO true "product resquest"
// @Success 201
// @Failure 400,500 {object} Error
// @Router /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductDTO
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	
	}
	p,err := entity.NewProduct(product.Name,product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) getProductByID(idStr string) (*entity.Product, error) {
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return nil, err
    }

    product, err := h.ProductDB.FindByID(uint(id))
    if err != nil {
		return nil,fmt.Errorf("product not found")
    }

    return product, nil
}

// Get a product by ID godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Success 200 {object} dto.CreateProductDTO
// @Failure 404
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	product, err := h.getProductByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update a product by ID godoc
// @Summary Update a product by ID
// @Description Update a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Param request body dto.CreateProductDTO true "product resquest"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err= h.getProductByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product.ID = uint(id)
	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}	
	w.WriteHeader(http.StatusOK)
}

// Delete a product by ID godoc
// @Summary Delete a product by ID
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "product id"
// @Success 200
// @Failure 400,404,500 {object} Error
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err= h.getProductByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	
}

// Get all products godoc
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param sort query string false "sort"
// @Success 200 {object} []dto.CreateProductDTO
// @Failure 400,500 {object} Error
// @Router /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt,err := strconv.Atoi(limit)
	if err != nil{
		limitInt = 10
	}
	
	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, limitInt,sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}
