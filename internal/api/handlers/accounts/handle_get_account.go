package accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (a *Accounts) HandleGetAcccount(ctx *gin.Context) {
	id := ctx.Param("id")
	account, err := a.store.GetAccount(ctx, id)

	if err != nil {
		log.Error().Err(err).Msg("failed to get account")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}

	ctx.JSON(http.StatusOK, newAccountResponse(account))
}
