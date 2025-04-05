package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, World")
	w.WriteHeader(http.StatusOK)
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {

	response := GetJsonProducts()
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func PostProductsHandler(w http.ResponseWriter, r *http.Request) {

	array, code := PostJsonProducts(r)
	if code == "201" {
		w.Write(array)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func Get_json_id(id int) []byte { //любой вывод json id
	var choto = map[string]int{"id": id}
	data, err := json.Marshal(choto)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func PostJsonProducts(r *http.Request) ([]byte, string) {

	var temp struct {
		Name   string
		Amount int
	}
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		return nil, "400"
	}

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO products (name, amount) VALUES ($1, $2)", temp.Name, temp.Amount)
	if err != nil {
		panic(err)
	}

	id, _ := result.LastInsertId()

	return Get_json_id(int(id)), "201"

}

func GetJsonProducts() []byte {

	var temp []struct {
		Id     int
		Name   string
		Amount int
	}

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", GetProductsHandler).Methods("GET")
	r.HandleFunc("/products", PostProductsHandler).Methods("POST")

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", r)
}
