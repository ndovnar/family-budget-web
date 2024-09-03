package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleDeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	claims := t.auth.GetClaimsFromContext(ctx)
	
	hasAccess := t.authz.IsUserHasWriteAcessToTransaction(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	transaction, err := t.store.GetTransaction(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("transaction not found")
		ctx.Error(error.NewHttpError(http.StatusNotFound))
		return
	}

	if transaction.UserID != claims.UserID {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	err = t.store.DeleteTransaction(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete transaction")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}
}
