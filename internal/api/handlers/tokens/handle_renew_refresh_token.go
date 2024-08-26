package tokens

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (t *Tokens) HandleRenewRefreshToken(ctx *gin.Context) {
	var req renewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	claims, err := t.auth.VerifyToken(req.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to verify token")
		ctx.Error(error.NewHttpError(http.StatusUnauthorized))
		return
	}

	session, err := t.store.GetSessionByID(ctx, claims.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	if session.IsDeleted {
		log.Error().Err(err).Msg("session is revoked")
		ctx.Error(error.NewHttpError(http.StatusUnauthorized))
		return
	}

	err = t.store.DeleteSession(ctx, session.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	newSession, err := t.store.CreateSession(ctx, &model.Session{
		UserID: session.UserID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	accessToken, err := t.auth.CreateAccessToken(newSession.ID, claims.UserID, claims.FirstName, claims.LastName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	refreshToken, err := t.auth.CreateRefreshToken(newSession.ID, claims.UserID, claims.FirstName, claims.LastName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create refresh token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	resp := newRenewRefreshTokenResponse(accessToken, refreshToken)
	ctx.JSON(http.StatusOK, resp)
}
