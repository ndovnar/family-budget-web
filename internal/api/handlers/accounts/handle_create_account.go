package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (a *Accounts) HandleCreateAccount(ctx *gin.Context) {
	claims := a.auth.GetClaimsFromContext(ctx)

	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	account, err := a.store.CreateAccount(ctx, &model.Account{
		Owner:   claims.UserID,
		Name:    req.Name,
		Balance: req.Balance,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create account")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newAccountResponse(account))
}
