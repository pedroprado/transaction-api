package entity

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAccountAddFunds(t *testing.T) {
	account := Account{Balance: 100}
	updated := account.AddFunds(200)

	expected := &Account{Balance: 300}

	assert.Equal(t, updated, expected)
}

func TestAccountRemoveFunds(t *testing.T) {
	account := Account{Balance: 100}

	updated, err := account.RemoveFunds(100)

	expected := &Account{Balance: 0}

	assert.Equal(t, nil, err)
	assert.Equal(t, updated, expected)
}

func TestAccountRemoveFunds_shouldNotRemoveFundWhenDoesNotHaveBalance(t *testing.T) {
	account := Account{Balance: 50}

	updated, err := account.RemoveFunds(100)

	var expected *Account

	assert.Equal(t, errorAccountDoesNotHaveFunds, err)
	assert.Equal(t, updated, expected)
}
