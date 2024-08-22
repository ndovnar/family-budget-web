package accounts

type createAccountRequest struct {
	Name    string  `json:"name" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

type updateAccountRequest struct {
	ID      string  `json:"id" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
	Owner   string  `json:"owner" binding:"required"`
}
