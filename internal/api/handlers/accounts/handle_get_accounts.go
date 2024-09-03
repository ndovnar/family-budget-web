package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/helpers/response"
	"github.com/rs/zerolog/log"
)

type GetAccountsQueryParams struct {
	Pagination filter.Pagination
}

func (a *Accounts) HandleGetAccounts(ctx *gin.Context) {
	filter := &filter.GetAccountsFilter{}
	if err := ctx.ShouldBindQuery(filter); err != nil {
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	claims := a.auth.GetClaimsFromContext(ctx)
	filter.OwnerID = claims.UserID

	accounts, count, err := a.store.GetAccounts(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get accounts")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	response.SetCountHeader(ctx, count)
	ctx.JSON(http.StatusOK, accounts)
}
