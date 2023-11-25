package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// instatiation of the server to run .
func (S *ApiServer) run() {
	router := mux.NewRouter()
	router.HandleFunc("/login", MakeHttpHandlerFunc(S.loginAccount))
	router.HandleFunc("/account", MakeHttpHandlerFunc(S.handleAccount))
	router.HandleFunc("/account/{id}", jwtMiddleWare(MakeHttpHandlerFunc(S.handleAccountsById), S.store))
	router.HandleFunc("/account/{id}/transfer", jwtMiddleWare(MakeHttpHandlerFunc(S.handleTransferAccount), S.store))
	log.Println("Port " + S.listenAddress + " is now running")
	http.ListenAndServe(S.listenAddress, router)

}
func (S *ApiServer) loginAccount(w http.ResponseWriter, req *http.Request) error {
	log := Login{}
	if err := json.NewDecoder(req.Body).Decode(&log); err != nil {
		return nil
	}
	acc, _ := S.store.GetAccountByNumber(int(log.Number))
	fmt.Println(acc)
	hash, err := createJWT(acc)
	if err != nil {
		permissionDenied(w)
	}
	if !acc.validPassword(hash) {
		permissionDenied(w)
	}
	ans := &LogInResponse{
		Token:  hash,
		Number: acc.Number,
	}
	token, _ := createJWT(acc)
	// if err!={}
	req.Header.Add("JWT-TOKEN", token)
	return writeJson(w, http.StatusOK, ans)

}
func jwtMiddleWare(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		validateTokenString := req.Header.Get("JWT-TOKEN")

		token, err := validateJwt(validateTokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userId, _ := getParams(req)
		account, err := s.GetAccountById(userId)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if account.ID != claims["accountNumber"] {
			permissionDenied(w)
			return
		}
		handlerFunc(w, req)
	}
}
func validateJwt(tokenVal string) (*jwt.Token, error) {
	token := os.Getenv("TOKEN")
	return jwt.Parse(tokenVal, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(token), nil
	})
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
}
func (S *ApiServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	values, err := S.store.GetsAllAccount()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, values)

}
func (S *ApiServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	id, err := getParams(req)
	if err != nil {
		return err
	}
	if _, err := S.store.GetAccountById(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, map[string]int{"deleted": id})

}
func (S *ApiServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {

	createAccount := new(CreateAccountRequest)

	if err := json.NewDecoder(req.Body).Decode(createAccount); err != nil {
		return err
	}
	account, err := HandleAccount(createAccount.FirstName, createAccount.LastName, createAccount.Password, createAccount.Email)
	if err != nil {

		log.Fatal(err)
	}
	if err := S.store.createAccount(account); err != nil {
		return err
	}
	// _, err = createJWT(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// req.Header.Add("JWT_TOKEN", tokenString)
	return writeJson(w, http.StatusOK, account)
}
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     2000000,
		"accountNumber": account.Number,
	}
	secret := os.Getenv("TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(secret))

}
func (S *ApiServer) handleTransferAccount(w http.ResponseWriter, req *http.Request) error {
	transfer := new(TransferRequest)
	if err := json.NewDecoder(req.Body).Decode(transfer); err != nil {
		return nil
	}
	S.store.FindAccount(transfer)
	defer req.Body.Close()
	return writeJson(w, http.StatusOK, transfer)
}

// To start the whole server
func main() {
	io.ReadAll()
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
