package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string
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

type ApiServer struct {
	listenAddress string
	store         Storage
}

func (S *ApiServer) run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandlerFunc(S.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(S.handleGetAccount))
	log.Println("Port " + S.listenAddress + " is now running")
	http.ListenAndServe(S.listenAddress, router)

}
func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
func newApiServer(las string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddress: las,
		store:         store,
	}
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
func (S *ApiServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	id := mux.Vars(req)
	fmt.Println(id)
	writeJson(w, http.StatusOK, &Account{})
	return nil
}
func (S *ApiServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func (S *ApiServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func main() {
	store, err := databaseConnection()
	log.Fatalf("error", err)
	// fmt.Println("Hello web server")
	w := newApiServer(":2000", store)
	w.run()
}
