package tokens

type renewAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func newRenewAccessTokenResponse(accessToken string) *renewAccessTokenResponse {
	return &renewAccessTokenResponse{
		AccessToken: accessToken,
	}
}

type renewRefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func newRenewRefreshTokenResponse(accessToken string, refreshToken string) *renewRefreshTokenResponse {
	return &renewRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
