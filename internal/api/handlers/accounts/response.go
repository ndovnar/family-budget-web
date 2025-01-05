package accounts

import "github.com/ndovnar/family-budget-api/internal/model"

type accountsResponse struct {
	Values []*model.Account
	Meta   *meta
}

type meta struct {
	Count int64 `json:"count"`
}

func newAccountsResponse(accounts []*model.Account, count int64) *accountsResponse {
	return &accountsResponse{
		Values: accounts,
		Meta: &meta{
			Count: count,
		},
	}
}
