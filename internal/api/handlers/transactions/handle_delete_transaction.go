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

	err := t.store.DeleteTransaction(ctx, id)
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
