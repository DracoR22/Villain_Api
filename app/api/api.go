package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/DracoR22/villain_api/app/types"
	"github.com/DracoR22/villain_api/storage"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// ----------------------------------------------------//ROUTES//---------------------------------------------//
func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Accounts
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/:id", makeHttpHandleFunc(s.handleGetAccountByID))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// ---------------------------------------------------//ACCOUNT ROUTES//-------------------------------------//
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not alowed %s", r.Method)
}

// -----------------------------------------------------//GET ACCOUNTS//--------------------------------------//
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// ---------------------------------------------------//GET ACCOUNT BY ID//------------------------------------//
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &types.Account{})
}

// ---------------------------------------------------//CREATE ACCOUNT//-------------------------------------//
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	CreateAccountReq := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(CreateAccountReq); err != nil {
		return err
	}

	// Check if FirstName and LastName are provided
	if CreateAccountReq.FirstName == "" || CreateAccountReq.LastName == "" {
		return errors.New("both FirstName and LastName are required")
	}

	account := types.NewAccount(CreateAccountReq.FirstName, CreateAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// ---------------------------------------------------//DELETE ACCOUNT//------------------------------------//
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// --------------------------------------------------//TURN DATA INTO JSON//------------------------------------//
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode((v))
}

// --------------------------------------------------//SETUP SERVER//--------------------------------------------//
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle The Error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
