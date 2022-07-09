package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)



type Product struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Rating int `json:"rating"`
}

var Products = make(map[int]Product)
var Id int = 0

func (p *Product) Modify(updated Product) {
	// log.Println("inside modify function")
	p.Name = updated.Name
	p.Price = updated.Price
	p.Rating = updated.Rating
	// log.Println(p)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Ecommerce")
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Products)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("creating a product")
	var product Product
	Id++
	product.Id = Id
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
	}
	Products[Id] = product
}

func getProductWithIdhandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if key, ok := vars["id"]; ok {
		id,err := strconv.Atoi(key)
		if err != nil {
			http.Error(w, "Error converting id in path to id of products", http.StatusBadRequest)
		}
		if product, ok := Products[id]; ok {
			json.NewEncoder(w).Encode(product)
		} else {
			http.Error(w, "Unable to find the product with given key", http.StatusBadRequest)
		}
	}
}

func updateProductWithId(w http.ResponseWriter, r *http.Request) {
	// log.Println("inside update function")
	vars := mux.Vars(r)
	var productToUpd Product
	if key, ok := vars["id"]; ok {
		id,err := strconv.Atoi(key)
		if err != nil {
			http.Error(w, "Error converting id in path to id of products", http.StatusBadRequest)
		}
		if _, ok := Products[id]; ok {
			productToUpd.Id = id
			err := json.NewDecoder(r.Body).Decode(&productToUpd)
			if err != nil {
				http.Error(w, "Error decoding the request to update", http.StatusBadRequest)
			}
			delete(Products, id)
			Products[id] = productToUpd
			// log.Println(Products)
			
		} else {
			http.Error(w, "Unable to find the product with given key", http.StatusBadRequest)
		}
	}
}

func deleteProductWithId(w http.ResponseWriter, r *http.Request) {
	// log.Println("inside update function")
	vars := mux.Vars(r)
	if key, ok := vars["id"]; ok {
		id,err := strconv.Atoi(key)
		if err != nil {
			http.Error(w, "Error converting id in path to id of products", http.StatusBadRequest)
		}
		if _ , ok := Products[id]; ok {
			delete(Products, id)
		} else {
			http.Error(w, "Unable to find the product with given key", http.StatusBadRequest)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/products", getProductsHandler).Methods("GET")
	r.HandleFunc("/create", createProductHandler).Methods("POST")
	r.HandleFunc("/products/{id}", getProductWithIdhandler).Methods("GET")
	r.HandleFunc("/products/{id}", updateProductWithId).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProductWithId).Methods("DELETE")


	server := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	log.Println("Listening to server...")
	log.Fatal(server.ListenAndServe())
}