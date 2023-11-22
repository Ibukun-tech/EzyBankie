package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandlerFunc(A ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := A(w, req); err != nil {
			writeJson(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}
	}
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
func newApiServer(las string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddress: las,
		store:         store,
	}
}

type ApiServer struct {
	listenAddress string
	store         Storage
}

func (S *ApiServer) run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandlerFunc(S.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(S.handleAccountsById))
	router.HandleFunc("/account/transfer", makeHttpHandlerFunc(S.handleTransferAccount))
	log.Println("Port " + S.listenAddress + " is now running")
	http.ListenAndServe(S.listenAddress, router)

}

func (S *ApiServer) handleAccountsById(w http.ResponseWriter, req *http.Request) error {
	if req.Method == "GET" {
		return S.handleGetAccountId(w, req)
	}
	if req.Method == "DELETE" {
		return S.handleDeleteAccount(w, req)
	}
	return nil
}

// Handling of request
func (S *ApiServer) handleAccount(w http.ResponseWriter, req *http.Request) error {
	if req.Method == "GET" {
		return S.handleGetAccount(w, req)
	}
	if req.Method == "POST" {
		return S.handleCreateAccount(w, req)
	}
	if req.Method == "DELETE" {
		return S.handleDeleteAccount(w, req)
	}
	return nil
}
func getParams(req *http.Request) (int, error) {
	vars := mux.Vars(req)["id"]
	id, err := strconv.Atoi(vars)
	return id, err
}
func (S *ApiServer) handleGetAccountId(w http.ResponseWriter, req *http.Request) error {

	id, err := getParams(req)
	if err != nil {
		return fmt.Errorf("invalid ID %d", id)
	}
	account, err := S.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, account)
	return nil
}
func (S *ApiServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	values, err := S.store.GetsAllAccount()
	if err != nil {
		return err
	}
	writeJson(w, http.StatusOK, values)
	return nil
}
func (S *ApiServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	id, err := getParams(req)
	if err != nil {
		return err
	}
	if _, err := S.store.GetAccountById(id); err != nil {
		return err
	}
	writeJson(w, http.StatusOK, map[string]int{"deleted": id})
	return nil
}
func (S *ApiServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {

	createAccount := new(CreateAccountRequest)

	if err := json.NewDecoder(req.Body).Decode(createAccount); err != nil {
		return err
	}
	account := HandleAccount(createAccount.FirstName, createAccount.LastName)
	if err := S.store.createAccount(account); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, account)
}

func (S *ApiServer) handleTransferAccount(w http.ResponseWriter, req *http.Request) error {
	transfer := new(TransferRequest)
	if err := json.NewDecoder(req.Body).Decode(transfer); err != nil {
		return nil
	}
	defer req.Body.Close()
	return writeJson(w, http.StatusOK, transfer)
}
func main() {
	store, err := databaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	if err = store.createAccoutTable(); err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Hello web server")
	w := newApiServer(":2000", store)
	w.run()
}
