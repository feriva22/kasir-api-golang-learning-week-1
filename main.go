package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Produk struct {
	ID    int     `json:"id"`
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stock"`
}

var produk = []Produk{
	{ID: 1, Nama: "Obeng", Harga: 50000, Stok: 100},
	{ID: 2, Nama: "Mur", Harga: 100, Stok: 1000},
}

// GET /api/produk/{id}
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk Belum Ada", http.StatusBadRequest)
}

// PUT /api/produk/{id}
func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//loop produk cari yang id nya sama
	for i := range produk {
		if produk[i].ID == id {
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Produk Belum Ada", http.StatusBadRequest)
}

// DELETE /api/produk/{id}
func deleteProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//loop produk cari yang id nya sama
	for i, p := range produk {
		if p.ID == id {

			//buat slice list produk sebelum dihapus produknya dan setelahnya
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Berhasil delete id",
			})
			return
		}
	}

	http.Error(w, "Produk Belum Ada", http.StatusBadRequest)
}

func sendEmail(to string, subject string, body string) error {
	fmt.Printf("[GOROUTINE STARTED] Sending email to: %s\nSubject: %s\nBody: %s\n", to, subject, body)

	// Simulate email sending delay
	time.Sleep(2 * time.Second)

	fmt.Printf("[GOROUTINE COMPLETED] Email sent to: %s\n", to)
	return nil
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// get detail produk
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getProdukByID(w, r)
		}

		if r.Method == "PUT" {
			updateProdukByID(w, r)
		}

		if r.Method == "DELETE" {
			deleteProdukByID(w, r)
		}
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			//baca dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			//masukkan data ke var
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			//kirim notif jika sudah tambah produk dengan goroutine
			fmt.Printf("[MAIN] Starting goroutine to send email for product: %s\n", produkBaru.Nama)
			go sendEmail("user@gmail.com", "Sukses tambah produk baru "+produkBaru.Nama, "Kamu sukses menambah produk")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //201

			json.NewEncoder(w).Encode(produk)
		}

	})

	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
