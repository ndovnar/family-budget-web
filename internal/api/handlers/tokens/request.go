package tokens

type renewTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
