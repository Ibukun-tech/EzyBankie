package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCreatAccount(t *testing.T) {
	acc, err := HandleAccount("s", "asd", "ahshss", "hshsh")
	assert.Nil(t, err)
	fmt.Println(acc)
}
