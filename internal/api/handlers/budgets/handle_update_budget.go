package budgets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (b *Budgets) HandleUpdateBudget(ctx *gin.Context) {
	id := ctx.Param("id")

	hasAccess := b.authz.IsUserHasAccessToBudget(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	var req budgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	budget, err := b.store.UpdateBudget(ctx, id, &model.Budget{
		Name: req.Name,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to updated budget")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, budget)
}
