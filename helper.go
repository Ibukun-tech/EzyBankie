package main

import (
	"encoding/json"
	"net/http"
)

func permissionDenied(w http.ResponseWriter) {
	writeJson(w, http.StatusBadRequest, ApiError{
		Error: "Invalid token",
	})
}

type ApiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string `json:"error"`
}

func MakeHttpHandlerFunc(A ApiFunc) http.HandlerFunc {
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
