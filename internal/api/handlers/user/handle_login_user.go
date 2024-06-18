package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/hash"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (user *User) HandleLoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("login user: failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	savedUser, err := user.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("login user: failed to get user from DB")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
			return
		}

		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	err = hash.CheckPassword(req.Password, savedUser.Password)
	if err != nil {
		ctx.Error(error.NewHttpError(http.StatusUnauthorized))
		return
	}

	session, err := user.store.CreateSession(ctx, &model.Session{
		UserID: savedUser.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("login user: failed to create session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	accessToken, err := user.auth.CreateAccessToken(session.ID, savedUser.ID, savedUser.FirstName, savedUser.LastName)
	if err != nil {
		log.Error().Err(err).Msg("login user: failed to create access token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	refreshToken, err := user.auth.CreateRefreshToken(session.ID, savedUser.ID, savedUser.FirstName, savedUser.LastName)
	if err != nil {
		log.Error().Err(err).Msg("login user: failed to create refresh token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	resp := newLoginResponse(accessToken, refreshToken)
	ctx.JSON(http.StatusOK, resp)
}
