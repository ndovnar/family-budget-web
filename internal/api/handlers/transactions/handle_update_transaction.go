package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleUpdateTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	var req transactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	transaction, err := t.store.UpdateTransaction(ctx, id, &model.Transaction{
		Type:        req.Type,
		FromAccount: req.FromAccount,
		ToAccount:   req.ToAccount,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to updated transaction")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
