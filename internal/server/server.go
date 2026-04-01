package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-ticker/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	port   int
	limits Limits
}

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.mux.HandleFunc("POST /api/accounts", s.hCreateAcct)
	s.mux.HandleFunc("GET /api/accounts", s.hListAccts)
	s.mux.HandleFunc("DELETE /api/accounts/{id}", s.hDelAcct)
	s.mux.HandleFunc("POST /api/transactions", s.hAddTxn)
	s.mux.HandleFunc("GET /api/transactions", s.hListTxns)
	s.mux.HandleFunc("DELETE /api/transactions/{id}", s.hDelTxn)
	s.mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) })
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]string{"status": "ok"}) })
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"product": "stockyard-ticker", "version": "0.1.0"})
	})
	return s
}

func (s *Server) Start() error {
	log.Printf("[ticker] :%d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) hCreateAcct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Currency string `json:"currency"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Name == "" {
		wj(w, 400, map[string]string{"error": "name required"})
		return
	}
	a, err := s.db.CreateAccount(req.Name, req.Type, req.Currency)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"account": a})
}

func (s *Server) hListAccts(w http.ResponseWriter, r *http.Request) {
	a, _ := s.db.ListAccounts()
	if a == nil {
		a = []store.Account{}
	}
	wj(w, 200, map[string]any{"accounts": a, "count": len(a)})
}

func (s *Server) hDelAcct(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteAccount(r.PathValue("id"))
	wj(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) hAddTxn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID   string `json:"account_id"`
		Description string `json:"description"`
		AmountCents int    `json:"amount_cents"`
		Category    string `json:"category"`
		Date        string `json:"date"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.AccountID == "" || req.Description == "" {
		wj(w, 400, map[string]string{"error": "account_id and description required"})
		return
	}
	t, err := s.db.AddTransaction(req.AccountID, req.Description, req.AmountCents, req.Category, req.Date)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"transaction": t})
}

func (s *Server) hListTxns(w http.ResponseWriter, r *http.Request) {
	t, _ := s.db.ListTransactions(r.URL.Query().Get("account_id"), 50)
	if t == nil {
		t = []store.Transaction{}
	}
	wj(w, 200, map[string]any{"transactions": t, "count": len(t)})
}

func (s *Server) hDelTxn(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteTransaction(r.PathValue("id"))
	wj(w, 200, map[string]string{"status": "deleted"})
}

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
