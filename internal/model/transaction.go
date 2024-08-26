package model

type Transaction struct {
	ID          string          `json:"id,omitempty" bson:"_id,omitempty"`
	Type        TransactionType `json:"type" bson:"type"`
	FromAccount string          `json:"fromAccount" bson:"fromAccount"`
	ToAccount   string          `json:"toAccount" bson:"toAccount"`
	Category    string          `json:"category" bson:"category"`
	User        string          `json:"user" bson:"user"`
	Amount      float64         `json:"amount" bson:"amount"`
	Description string          `json:"description" bson:"description"`
	IsDeleted   bool            `json:"deleted,omitempty" bson:"deleted"`
	Dates       Dates           `json:"dates" bson:"dates"`
}

type TransactionType string

const (
	TransactionTypeTransfer = TransactionType("transfer")
	TransactionTypeExpense  = TransactionType("expense")
	TransactionTypeIncome   = TransactionType("income")
)

type GetTransactionsFilter struct {
	Account string
	Deleted bool
}
