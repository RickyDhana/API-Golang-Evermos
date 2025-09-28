package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"project_satu/internal/config"
	"project_satu/internal/domain"
	"project_satu/internal/handler"
	"project_satu/internal/middleware"
)

func main() {
	// load environment variables
	_ = godotenv.Load()

	// init database
	config.InitDB()
	defer config.DB.Close()

	// auto migrate semua model
	err := config.DB.AutoMigrate(
		&domain.User{},
		&domain.Toko{},
		&domain.Alamat{},
		&domain.Category{},
		&domain.Produk{},
		&domain.FotoProduk{},
		&domain.Trx{},
		&domain.DetailTrx{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.ProductLog{},
	).Error
	if err != nil {
		log.Fatal("AutoMigrate error: ", err)
	}

	// setup router
	r := mux.NewRouter()

	// default root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Mini Evermos API Running...")
	}).Methods("GET")

	// public routes
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	// protected routes (pakai middleware JWT)
	r.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(handler.MeHandler))).Methods("GET")
	r.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateMeHandler))).Methods("PUT")

	// alamat routes (protected)
	r.Handle("/alamat", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateAlamatHandler))).Methods("POST")
	r.Handle("/alamat", middleware.AuthMiddleware(http.HandlerFunc(handler.ListAlamatHandler))).Methods("GET")
	r.Handle("/alamat/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateAlamatHandler))).Methods("PUT")
	r.Handle("/alamat/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteAlamatHandler))).Methods("DELETE")

	// toko routes (protected)
	r.Handle("/toko", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateTokoHandler))).Methods("POST")
	r.Handle("/toko", middleware.AuthMiddleware(http.HandlerFunc(handler.ListTokoHandler))).Methods("GET")
	r.Handle("/toko/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateTokoHandler))).Methods("PUT")
	r.Handle("/toko/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteTokoHandler))).Methods("DELETE")

	// kategori routes
	r.Handle("/kategori", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateKategoriHandler))).Methods("POST")
	r.HandleFunc("/kategori", handler.ListKategoriHandler).Methods("GET")
	r.Handle("/kategori/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateKategoriHandler))).Methods("PUT")
	r.Handle("/kategori/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteKategoriHandler))).Methods("DELETE")

	// product routes
	r.Handle("/product", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateProdukHandler))).Methods("POST")
	r.HandleFunc("/product", handler.ListProdukHandler).Methods("GET") // public
	r.HandleFunc("/product/{id}", handler.GetProdukHandler).Methods("GET") // public
	r.Handle("/product/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateProdukHandler))).Methods("PUT")
	r.Handle("/product/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteProdukHandler))).Methods("DELETE")
	r.Handle("/product/{id}/upload", middleware.AuthMiddleware(http.HandlerFunc(handler.UploadProdukPhotoHandler))).Methods("POST")

	// order routes (protected)
r.Handle("/orders", middleware.AuthMiddleware(http.HandlerFunc(handler.CreateOrderHandler))).Methods("POST")
r.Handle("/orders", middleware.AuthMiddleware(http.HandlerFunc(handler.ListOrderHandler))).Methods("GET")
r.Handle("/orders/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetOrderHandler))).Methods("GET")


	// start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
