package tokens

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/rs/zerolog/log"
)

func (t *Tokens) HandleRenewAccessToken(ctx *gin.Context) {
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

	accessToken, err := t.auth.CreateAccessToken(session.ID, claims.UserID, claims.FirstName, claims.LastName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access token")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	resp := newRenewAccessTokenResponse(accessToken)
	ctx.JSON(http.StatusOK, resp)
}
