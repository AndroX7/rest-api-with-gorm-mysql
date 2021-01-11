package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"string"`
}

func main() {

	db, err = gorm.Open("mysql", "root:@/go_rest_api_crud?charset=utf8&parseTime=true")
	if err != nil {
		log.Println("Connection Failed!", err)
	} else {
		log.Println("Connection Successed")
	}

	db.AutoMigrate(&Product{})
	handleRequests()
}

func handleRequests() {
	log.Println("start development at http://localhost:3009")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", editProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3009", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Create Product")
	payloads, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var product Product
	err = json.Unmarshal(payloads, &product)

	if err != nil {
		log.Println(err)
	}

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success Get Products"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)

	res := Result{Code: 200, Data: product, Message: "Success Get Product"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	payloads, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var productUpdate Product
	err = json.Unmarshal(payloads, &productUpdate)

	if err != nil {
		log.Println(err)
	}

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdate)
	// db.Create(&productUpdate)

	res := Result{Code: 200, Data: product, Message: "Success Update Product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)
	res := db.Delete(&product)

	response := Result{Code: 200, Data: res, Message: "Success Delete Product"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
