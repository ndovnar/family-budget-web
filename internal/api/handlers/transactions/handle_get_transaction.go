package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleGetTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	hasAccess := t.authz.IsUserHasReadAcessToTransaction(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	transaction, err := t.store.GetTransaction(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get transaction")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
