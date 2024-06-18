package account

import "github.com/ndovnar/family-budget-api/internal/model"

type accountResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

func newAccountResponse(account *model.Account) *accountResponse {
	return &accountResponse{
		ID:      account.ID,
		Name:    account.Name,
		Balance: account.Balance,
	}
}
