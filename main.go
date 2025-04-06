package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, World")
	w.WriteHeader(http.StatusOK)
}

type Product struct {
	Name   string
	Amount int
}

type Controller struct {
	Service *Service
}

type Service struct {
	Repository *Repository
}

type Repository struct {
	DB *sql.DB
}

func NewController(service *Service) *Controller {
	return &Controller{
		Service: service,
	}
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db
}



func (c Controller) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := c.Service.GetAllProducts()

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}





func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", GetProductsHandler).Methods("GET")

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", r)
}
