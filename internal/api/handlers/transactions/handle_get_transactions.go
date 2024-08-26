package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (t *Transactions) HandleGetTransactions(ctx *gin.Context) {
	account, _ := ctx.GetQuery("account")
	filter := &model.GetTransactionsFilter{Account: account}

	transactions, err := t.store.GetTransactions(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get transactions")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}
