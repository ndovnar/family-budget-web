package filter

type GetTransactionsFilter struct {
	FromAccountID string `form:"fromAccount"`
	ToAccountID   string `form:"toAccount"`
	AccountID     string `form:"account"`
	CategoryID    string `form:"category"`
	UserID        string
	Pagination    *Pagination
}
