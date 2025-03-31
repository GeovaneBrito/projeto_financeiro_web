package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Asset representa um ativo financeiro.
type Asset struct {
	Type       string  `json:"type"`
	Ticker     string  `json:"ticker"`
	Price      float64 `json:"price"`
	Percentage float64 `json:"percentage"`
	Score      int     `json:"score"`
	Quantity   int     `json:"quantity"`
}

// Portfolio representa os dados consolidados da carteira.
type Portfolio struct {
	TotalValue     float64            `json:"totalValue"`
	Distribution   map[string]float64 `json:"distribution"`
	AssetsQuantity int                `json:"assetsQuantity"`
}

// assets é uma lista mockada de ativos.
var assets = []Asset{
	{"Fundos Imobiliários", "KNSC11", 3664.00, 19.58, -2, 400},
	{"Ações Nacionais", "ITUB4", 2436.28, 13.04, 13, 77},
	{"Ações Nacionais", "USIM5", 2070.64, 11.08, 11, 362},
	{"Ações Internacionais", "KOF", 536.10, 4.35, 5, 1048},
	{"Ações Nacionais", "ABEV3", 1655.40, 8.86, 9, 124},
}

// getAssets retorna a lista de ativos.
func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

// getAssetByTicker busca um ativo pelo ticker.
func getAssetByTicker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticker := vars["ticker"]

	for _, asset := range assets {
		if asset.Ticker == ticker {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(asset)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ativo não encontrado"})
}

// addAsset adiciona um novo ativo.
func addAsset(w http.ResponseWriter, r *http.Request) {
	var newAsset Asset
	err := json.NewDecoder(r.Body).Decode(&newAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assets = append(assets, newAsset)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAsset)
}

// editAsset edita um ativo existente.
func editAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticker := vars["ticker"]

	var updatedAsset Asset
	err := json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, asset := range assets {
		if asset.Ticker == ticker {
			assets[i] = updatedAsset
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedAsset)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ativo não encontrado"})
}

// getPortfolio retorna os dados consolidados da carteira.
func getPortfolio(w http.ResponseWriter, r *http.Request) {
	totalValue := 0.0
	distribution := make(map[string]float64)
	assetsQuantity := len(assets)

	for _, asset := range assets {
		totalValue += asset.Price * float64(asset.Quantity)
		distribution[asset.Type] += asset.Price * float64(asset.Quantity)
	}

	for assetType, value := range distribution {
		distribution[assetType] = (value / totalValue) * 100
	}

	portfolio := Portfolio{
		TotalValue:     totalValue,
		Distribution:   distribution,
		AssetsQuantity: assetsQuantity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portfolio)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/assets", getAssets).Methods("GET")
	r.HandleFunc("/api/assets/{ticker}", getAssetByTicker).Methods("GET")
	r.HandleFunc("/api/assets/add", addAsset).Methods("POST")
	r.HandleFunc("/api/assets/edit/{ticker}", editAsset).Methods("PUT")
	r.HandleFunc("/api/portfolio", getPortfolio).Methods("GET")

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
