package accounts

import (
	"github.com/ndovnar/family-budget-api/internal/model"
)

type accountResponse struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Dates   model.Dates `json:"dates"`
	Balance float64     `json:"balance"`
}

func newAccountResponse(account *model.Account) *accountResponse {
	return &accountResponse{
		ID:      account.ID,
		Name:    account.Name,
		Balance: account.Balance,
		Dates:   account.Dates,
	}
}

func newAccountsResponse(accounts []*model.Account) []*accountResponse {
	accountsResponse := make([]*accountResponse, 0, len(accounts))
	for _, account := range accounts {
		innerAccount := account
		accountsResponse = append(accountsResponse, newAccountResponse(innerAccount))
	}

	return accountsResponse
}
