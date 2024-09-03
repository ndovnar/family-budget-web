package filter

type GetTransactionsFilter struct {
	FromAccountID string `form:"fromAccount"`
	ToAccountID   string `form:"toAccount"`
	CategoryID    string `form:"category"`
	UserID        string
	Pagination    *Pagination
}
