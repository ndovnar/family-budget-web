package transactions

import "github.com/ndovnar/family-budget-api/internal/model"

const ColName = "name”"

type transactionRequest struct {
	Type        model.TransactionType `json:"type" binding:"required,oneof=transfer expense income"`
	FromAccount string                `json:"fromAccount" binding:"required"`
	ToAccount   string                `json:"toAccount" binding:"required_if=Type transfer"`
	Category    string                `json:"category"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
}
