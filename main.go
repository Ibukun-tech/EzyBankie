package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error

func handleError(A ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := A(w, req); err != nil {
			// Handle the error Ibk
		}
	}
}

type ApiServer struct {
	listenAddress string
}

func (S *ApiServer) run() {
	router := mux.NewRouter()
	router.HandleFunc("/account")
}
func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
func newApiServer(las string) *ApiServer {
	return &ApiServer{
		listenAddress: las,
	}
}

// Handling of request
func (S *ApiServer) handleAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func (S *ApiServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func (S *ApiServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func (S *ApiServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}
func main() {
	fmt.Println("Hello web server")
}
