package users

import "github.com/ndovnar/family-budget-api/internal/model"

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func newUserResponse(user *model.User) *userResponse {
	return &userResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func newLoginResponse(accessToken string, refreshToken string) *loginResponse {
	return &loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
