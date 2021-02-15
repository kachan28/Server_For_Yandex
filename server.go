package main

import (
	"authorization"
	"fmt"
	"log"
	"net/http"
	"products"
	"registration"

	"github.com/gorilla/mux"
)

func notfoundhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method not found")
}

func main() {
	r := mux.NewRouter().StrictSlash(false)

	regHandler := http.HandlerFunc(registration.RegistrationHandler)
	authHandler := http.HandlerFunc(authorization.AuthHandler)
	refreshHandler := http.HandlerFunc(authorization.RefreshToken)
	productsInfoHandler := http.HandlerFunc(products.ProductInfoHandler)
	mainPageHandler := http.HandlerFunc(products.GetProductsForMainHandler)
	findHandler := http.HandlerFunc(products.FindProductHandler)

	r.HandleFunc("/registration", regHandler).Methods("POST")
	r.HandleFunc("/login", authHandler).Methods("POST")
	r.HandleFunc("/refresh", refreshHandler).Methods("POST")
	productsPath := r.PathPrefix("/products").Subrouter()
	productsPath.Handle("/product/{shop}/{articul}", authorization.JWTAccessHandler(productsInfoHandler)).Methods("GET")
	productsPath.Handle("/main", authorization.JWTAccessHandler(mainPageHandler)).Methods("GET")
	productsPath.Handle("/find", authorization.JWTAccessHandler(findHandler)).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(notfoundhandler).GetHandler()

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}
	log.Println("listening...")
	server.ListenAndServe()
}
