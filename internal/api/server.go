package api

import (
	"encoding/json"
	"modular-task/internal/products"
	"modular-task/internal/users"
	"net/http"
	"strconv"
)

type Server struct {
	userService    *users.UserService
	productService *products.ProductService
}

func NewServer(u *users.UserService, p *products.ProductService) *Server {
	return &Server{
		userService:    u,
		productService: p,
	}
}

func (s *Server) Routes() {
	http.HandleFunc("/users", s.handleUsers)
	http.HandleFunc("/products", s.handleProducts)
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.userService.CreateUser(input.Name, input.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)

	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		user, ok := s.userService.GetUserByID(int64(id))
		if !ok {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			UserID int64   `json:"user_id"`
			Name   string  `json:"name"`
			Price  float64 `json:"price"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := s.productService.CreateProduct(input.UserID, input.Name, input.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(product)
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		product, ok := s.productService.GetProductByID(int64(id))
		if !ok {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(product)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
