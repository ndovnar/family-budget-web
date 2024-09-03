package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (u *Users) HandleLoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	user, err := u.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user from DB")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
			return
		}

		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	err = checkPassword(req.Password, user.Password)
	if err != nil {
		log.Error().Err(err).Msg("password is incorrect")
		ctx.Error(error.NewHttpError(http.StatusUnauthorized))
		return
	}

	session, err := u.store.CreateSession(ctx, &model.Session{
		UserID: user.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	accessToken, err := u.auth.CreateAccessToken(session.ID, user.ID, user.FirstName, user.LastName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	refreshToken, err := u.auth.CreateRefreshToken(session.ID, user.ID, user.FirstName, user.LastName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create refresh token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	resp := newLoginResponse(accessToken, refreshToken)
	ctx.JSON(http.StatusOK, resp)
}
