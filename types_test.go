package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCreatAccount(t *testing.T) {
	acc, err := HandleAccount("s", "asd", "ahshss", "hshsh")
	assert.Nil(t, err)
	fmt.Println(acc)
}

func TestNewApiServer(t *testing.T) {
	store, err := databaseConnection()
	w := newApiServer(":3000", store)
	assert.Nil(t, err)
	fmt.Println(w)
}

func TestMakeHttpHandlerFunc(t *testing.T) {
	w := func(w http.ResponseWriter, req *http.Request) error {
		return nil
	}
	MakeHttpHandlerFunc(w)
}

func TestWriteJsonFunc(t *testing.T) {
	type value struct {
		Number int `json:"nm"`
	}
	type res http.ResponseWriter
	name := &value{
		Number: 2,
	}
	var b res
	writeJson(b, http.StatusAccepted, name)
	t.Helper()
}
