package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleCreateTransaction(ctx *gin.Context) {
	claims := t.auth.GetClaimsFromContext(ctx)

	var req transactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	transaction, err := t.store.CreateTransaction(ctx, &model.Transaction{
		Type:          req.Type,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		CategoryID:    req.CategoryID,
		UserID:        claims.UserID,
		Amount:        req.Amount,
		Description:   req.Description,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create transaction")
		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusBadRequest))
		}

		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
