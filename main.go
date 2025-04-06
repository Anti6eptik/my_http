package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
	_ "modernc.org/sqlite"
)

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
	return db
}

func (s Service) GetAllProducts() []byte {
	data := s.Repository.GetAllProducts()
	array, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return array

}

func (r Repository) GetAllProducts() []struct {
	Id     int
	Name   string
	Amount int
} {
	rows, err := r.DB.Query("select * from products")
	if err != nil {
		panic(err)
	}

	var Temp struct {
		Id     int
		Name   string
		Amount int
	}

	var result []struct {
		Id     int
		Name   string
		Amount int
	}

	for rows.Next() {
		rows.Scan(&Temp.Id, &Temp.Name, &Temp.Amount)
		result = append(result, Temp)
	}

	return result
}

func (c Controller) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	response := c.Service.GetAllProducts()
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
