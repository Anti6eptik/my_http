package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
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

func (s Service) GetAllProducts() {
	response := s.Repository.GetAllProducts
}

func (r Repository) GetAllProducts() []struct {
	id     int
	name   string
	amount int
} {
	rows, err := r.DB.Query("select * from products")
	if err != nil {
		panic(err)
	}

	var temp struct {
		id     int
		name   string
		amount int
	}

	var result []struct {
		id     int
		name   string
		amount int
	}

	for rows.Next() {
		rows.Scan(&temp.id, &temp.name, &temp.amount)
		result = append(result, temp)
	}

	return result
}

func (c Controller) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := c.Service.GetAllProducts()
	if err != nil {
		panic(err)
	}
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func main() {
	container := dig.New()
	_ = container.Provide(NewController)
	_ = container.Provide(NewService)
	_ = container.Provide(NewRepository)
	_ = container.Provide(NewDB)

	container.Invoke(func(controller *Controller) {
		r := mux.NewRouter()
		r.HandleFunc("/products", controller.GetProductsHandler).Methods("GET")
		fmt.Println("Server is listening...")
		http.ListenAndServe(":8080", r)
	})
}
