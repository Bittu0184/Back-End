package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type product struct {
	Product_id          string `json:"product_id"`
	Product_name        string `json:"product_name"`
	Image_path          string `json:"image_path"`
	Product_price       int    `json:"product_price"`
	Available           int    `json:"available"`
	Product_description string `json:"product_description"`
	Product_category    string `json:"product_category"`
}

type category struct {
	Category_name        string `json:"category_name"`
	Category_Id          int    `json:"category_id"`
	Category_description string `json:"category_description"`
	Image_path           string `json:"image_path"`
}

var db *sql.DB

func main() {
	fmt.Print("Hey There!")
	myrouter := mux.NewRouter().StrictSlash(true)
	db, _ = sql.Open("mysql", "")
	myrouter.HandleFunc("/get_product_category", GetProduct).Methods("POST")
	myrouter.HandleFunc("/get_category", GetCategory)
	http.ListenAndServe(":8081", myrouter)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	res, _ := ioutil.ReadAll(r.Body)
	var category string
	json.Unmarshal(res, &category)
	prod := GetProductWithCategory(category)
	json.NewEncoder(w).Encode(prod)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	AllCategory := GetAllCategory()
	json.NewEncoder(w).Encode(AllCategory)

}

func GetAllCategory() []category {
	rows, err := db.Query("SELECT * FROM product_category")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var categorys []category
	for rows.Next() {
		var cat category
		if err := rows.Scan(&cat.Category_name, &cat.Category_Id, &cat.Category_description, &cat.Image_path); err != nil {
			log.Fatal(err)
		}
		categorys = append(categorys, cat)
	}
	return categorys
}

func GetProductWithCategory(category string) []product {
	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var products []product
	for rows.Next() {
		var prod product
		if err := rows.Scan(&prod.Product_id, &prod.Product_name, &prod.Image_path, &prod.Product_price, &prod.Available, &prod.Product_description, &prod.Product_category); err != nil {
			log.Fatal(err)
		}
		if prod.Product_category == category {
			products = append(products, prod)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return products
}
