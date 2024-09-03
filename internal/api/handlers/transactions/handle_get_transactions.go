package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/helpers/response"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleGetTransactions(ctx *gin.Context) {
	filter := &filter.GetTransactionsFilter{}
	if err := ctx.ShouldBindQuery(filter); err != nil {
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	if filter.CategoryID == "" && filter.FromAccountID == "" && filter.ToAccountID == "" {
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	if filter.FromAccountID != "" {
		hasAccess := t.authz.IsUserHasAccessToAccount(ctx, filter.FromAccountID)

		if !hasAccess {
			ctx.Error(error.NewHttpError(http.StatusForbidden))
			return
		}
	}

	if filter.ToAccountID != "" {
		hasAccess := t.authz.IsUserHasAccessToAccount(ctx, filter.ToAccountID)

		if !hasAccess {
			ctx.Error(error.NewHttpError(http.StatusForbidden))
			return
		}
	}

	if filter.CategoryID != "" {
		hasAccess := t.authz.IsUserHasAccessToAccount(ctx, filter.CategoryID)

		if !hasAccess {
			ctx.Error(error.NewHttpError(http.StatusForbidden))
			return
		}
	}

	claims := t.auth.GetClaimsFromContext(ctx)
	filter.UserID = claims.UserID

	transactions, count, err := t.store.GetTransactions(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get transactions")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	response.SetCountHeader(ctx, count)
	ctx.JSON(http.StatusOK, transactions)
}
