package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/rs/zerolog/log"
)

func (u *Users) HandleLogoutUser(ctx *gin.Context) {
	claims := u.auth.GetClaimsFromContext(ctx)

	err := u.store.DeleteSession(ctx, claims.SessionID)
	if err != nil {
		log.Debug().Msg("failed to delete session")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.Status(http.StatusOK)
}
