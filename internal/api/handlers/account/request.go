package account

type newAccountRequest struct {
	Name    string  `json:"name" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}
