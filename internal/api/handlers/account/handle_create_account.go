package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (account *Account) HandleCreateAccount(ctx *gin.Context) {
	log.Debug().Msg("create account: handle")

	var req newAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("create account: failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	newAccount, err := account.store.CreateAccount(ctx, &model.Account{
		Name:    req.Name,
		Balance: req.Balance,
	})
	if err != nil {
		log.Error().Err(err).Msg("create account: failed to create")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	log.Debug().Msg("create account: success")
	ctx.JSON(http.StatusOK, newAccountResponse(newAccount))
}
