package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (a *Accounts) HandleUpdateAccount(ctx *gin.Context) {
	id := ctx.Param("id")

	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	account, err := a.store.UpdateAccount(ctx, id, &model.Account{
		Owner:   req.Owner,
		Name:    req.Name,
		Balance: req.Balance,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to updated account")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newAccountResponse(account))
}
