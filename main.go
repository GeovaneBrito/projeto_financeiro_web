package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors" // Importe o pacote cors
)

// Asset representa um ativo financeiro.
type Asset struct {
	ID         int     `json:"id"`
	Type       string  `json:"type"`
	Ticker     string  `json:"ticker"`
	Price      float64 `json:"price"`
	Percentage float64 `json:"percentage"`
	Score      int     `json:"score"`
	Quantity   int     `json:"quantity"`
}

// Contribution representa um aporte do usuário.
type Contribution struct {
	UserID int     `json:"userID"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

// InvestmentSuggestion representa uma sugestão de investimento.
type InvestmentSuggestion struct {
	AssetID           int     `json:"assetID"`
	SuggestedAmount   float64 `json:"suggestedAmount"`
	SuggestedQuantity int     `json:"suggestedQuantity"`
}

// Goal representa uma meta do usuário.
type Goal struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userID"`
	AssetType  string  `json:"assetType"`
	Percentage float64 `json:"percentage"`
}

// Question representa uma pergunta sobre um ativo.
type Question struct {
	ID        int    `json:"id"`
	Criterion string `json:"criterion"`
	Question  string `json:"question"`
	AssetType string `json:"assetType"`
}

// AssetQuestionAnswer representa a resposta de um usuário a uma pergunta sobre um ativo.
type AssetQuestionAnswer struct {
	AssetID    int    `json:"assetID"`
	QuestionID int    `json:"questionID"`
	Answer     string `json:"answer"`
}

// assets é uma lista mockada de ativos.
var assets = []Asset{
	{1, "Fundos Imobiliários", "XPML11", 100.95, 3.97, 6, 50},
	{2, "Fundos Imobiliários", "HFOF11", 54.50, 2.67, 4, 30},
	{3, "Ações Nacionais", "ITUB4", 31.64, 1.83, 13, 200},
}

// contributions é uma lista mockada de aportes.
var contributions = []Contribution{}

// goals é uma lista mockada de metas.
var goals = []Goal{}

// questions é uma lista mockada de perguntas.
var questions = []Question{}

// assetQuestionAnswers é uma lista mockada de respostas de usuários a perguntas sobre ativos.
var assetQuestionAnswers = []AssetQuestionAnswer{}

// getLatestContribution retorna o último aporte de um usuário.
func getLatestContribution(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var latestContribution Contribution
	for _, contribution := range contributions {
		if contribution.UserID == userID {
			latestContribution = contribution
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latestContribution)
}

// addContribution adiciona um novo aporte.
func addContribution(w http.ResponseWriter, r *http.Request) {
	var newContribution Contribution
	err := json.NewDecoder(r.Body).Decode(&newContribution)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contributions = append(contributions, newContribution)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newContribution)
}

// getInvestmentSuggestions retorna sugestões de investimento.
func getInvestmentSuggestions(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var suggestions []InvestmentSuggestion
	for _, asset := range assets {
		suggestedAmount := amount * (asset.Percentage / 100)
		suggestedQuantity := int(suggestedAmount / asset.Price)
		suggestions = append(suggestions, InvestmentSuggestion{
			AssetID:           asset.ID,
			SuggestedAmount:   suggestedAmount,
			SuggestedQuantity: suggestedQuantity,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// getGoals retorna as metas de um usuário.
func getGoals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userGoals []Goal
	for _, goal := range goals {
		if goal.UserID == userID {
			userGoals = append(userGoals, goal)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userGoals)
}

// addGoal adiciona uma nova meta.
func addGoal(w http.ResponseWriter, r *http.Request) {
	var newGoal Goal
	err := json.NewDecoder(r.Body).Decode(&newGoal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goals = append(goals, newGoal)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGoal)
}

// getQuestions retorna a lista de perguntas.
func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// addQuestion adiciona uma nova pergunta.
func addQuestion(w http.ResponseWriter, r *http.Request) {
	var newQuestion Question
	err := json.NewDecoder(r.Body).Decode(&newQuestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	questions = append(questions, newQuestion)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newQuestion)
}

// getAssetQuestionAnswers retorna as respostas de um usuário a perguntas sobre um ativo.
func getAssetQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID, err := strconv.Atoi(vars["assetID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var assetAnswers []AssetQuestionAnswer
	for _, answer := range assetQuestionAnswers {
		if answer.AssetID == assetID {
			assetAnswers = append(assetAnswers, answer)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assetAnswers)
}

// addAssetQuestionAnswer adiciona uma nova resposta de um usuário a uma pergunta sobre um ativo.
func addAssetQuestionAnswer(w http.ResponseWriter, r *http.Request) {
	var newAnswer AssetQuestionAnswer
	err := json.NewDecoder(r.Body).Decode(&newAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetQuestionAnswers = append(assetQuestionAnswers, newAnswer)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAnswer)
}

// getAssets retorna a lista de ativos, com opção de filtro por tipo.
func getAssets(w http.ResponseWriter, r *http.Request) {
	assetType := r.URL.Query().Get("type")

	var filteredAssets []Asset
	for _, asset := range assets {
		if assetType == "" || asset.Type == assetType {
			filteredAssets = append(filteredAssets, asset)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredAssets)
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
	assetID, err := strconv.Atoi(vars["assetID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedAsset Asset
	err = json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, asset := range assets {
		if asset.ID == assetID {
			assets[i] = updatedAsset
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedAsset)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ativo não encontrado"})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/contributions/latest/{userID}", getLatestContribution).Methods("GET")
	r.HandleFunc("/api/contributions", addContribution).Methods("POST")
	r.HandleFunc("/api/suggestions", getInvestmentSuggestions).Methods("GET")

	r.HandleFunc("/api/goals/{userID}", getGoals).Methods("GET")
	r.HandleFunc("/api/goals", addGoal).Methods("POST")

	r.HandleFunc("/api/questions", getQuestions).Methods("GET")
	r.HandleFunc("/api/questions", addQuestion).Methods("POST")

	r.HandleFunc("/api/asset-question-answers/{assetID}", getAssetQuestionAnswers).Methods("GET")
	r.HandleFunc("/api/asset-question-answers", addAssetQuestionAnswer).Methods("POST")

	r.HandleFunc("/api/assets", getAssets).Methods("GET")
	r.HandleFunc("/api/assets", addAsset).Methods("POST")
	r.HandleFunc("/api/assets/{assetID}", editAsset).Methods("PUT")

	// Configurar o CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Permite requisições de localhost:4200
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Permite os métodos HTTP necessários
		AllowedHeaders:   []string{"*"},                                       // Permite todos os cabeçalhos
	})
	// Usar o middleware CORS
	handler := c.Handler(r)

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
