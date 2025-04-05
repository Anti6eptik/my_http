package main

import (
	"database/sql"
	"encoding/json"
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

func GetProductsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	response := GetJsonProducts(db)
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func PostProductsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	array, code := PostJsonProducts(db, r)
	if code == "201" {
		w.Write(array)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func Get_json_id(id int) []byte {
	var choto = map[string]int{"id": id}
	data, err := json.Marshal(choto)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func PostJsonProducts(db *sql.DB, r *http.Request) ([]byte, string) {

	var temp struct {
		Name   string
		Amount int
	}
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		return nil, "400"
	}

	result, err := db.Exec("INSERT INTO products (name, amount) VALUES ($1, $2)", temp.Name, temp.Amount)
	if err != nil {
		panic(err)
	}

	id, _ := result.LastInsertId()

	return Get_json_id(int(id)), "201"

}

func GetJsonProducts(db *sql.DB) []byte {

	var temp []struct {
		Id     int
		Name   string
		Amount int
	}

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p struct {
			Id     int
			Name   string
			Amount int
		}
		err := rows.Scan(&p.Id, &p.Name, &p.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		temp = append(temp, p)
	}

	data, err := json.Marshal(temp)
	if err != nil {
		panic(err)
	}
	return data
}

func main() {
	r := mux.NewRouter()

	container := dig.New()

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	container.Provide(func() *sql.DB { return db })

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", container.invoke(GetProductsHandler)).Methods("GET")
	r.HandleFunc("/products", container.invoke(PostProductsHandler)).Methods("POST")

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", r)
}
