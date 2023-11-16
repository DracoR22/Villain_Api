package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DracoR22/villain_api/app/types"
	"github.com/DracoR22/villain_api/app/utils"
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
	router.HandleFunc("/account/{id}", WithJWTAuth(makeHttpHandleFunc(s.handleGetAccountByID)))

	// Transfer
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

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

	return fmt.Errorf("Method not alowed %s", r.Method)
}

// -----------------------------------------------------//GET ACCOUNTS//--------------------------------------//
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, accounts)
}

// ---------------------------------------------------//GET ACCOUNT BY ID//------------------------------------//
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {

		// Get Id from Params
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)

		if err != nil {
			return err
		}

		return utils.WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
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

	return utils.WriteJSON(w, http.StatusOK, account)
}

// ---------------------------------------------------//DELETE ACCOUNT//------------------------------------//
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Get Id from Params
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// -------------------------------------------------------//TRANSFER//---------------------------------------//
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return utils.WriteJSON(w, http.StatusOK, transferReq)
}

// -----------------------------------------------------//JWT//--------------------------------------------//

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling jwt middleware")

		tokenString := r.Header.Get("x-jwt-token")

		_, err := utils.ValidateJWT(tokenString)

		if err != nil {
			utils.WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle The Error
			utils.WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	// Get Id From Params
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}
