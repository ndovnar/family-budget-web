package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (a *Accounts) HandleGetAccounts(ctx *gin.Context) {
	claims := a.auth.GetClaimsFromContext(ctx)
	filter := &model.GetAccountsFilter{Owner: claims.UserID}

	accounts, err := a.store.GetAccounts(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get accounts")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
