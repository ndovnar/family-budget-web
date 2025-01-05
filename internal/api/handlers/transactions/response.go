package transactions

import "github.com/ndovnar/family-budget-api/internal/model"

type transactionsResponse struct {
	Values []*model.Transaction
	Meta   *meta
}

type meta struct {
	Count int64 `json:"count"`
}

func newTransactionsResponse(transactions []*model.Transaction, count int64) *transactionsResponse {
	return &transactionsResponse{
		Values: transactions,
		Meta: &meta{
			Count: count,
		},
	}
}
