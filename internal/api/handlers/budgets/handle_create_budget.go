package budgets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/rs/zerolog/log"
)

func (b *Budgets) HandleCreateBudget(ctx *gin.Context) {
	claims := b.auth.GetClaimsFromContext(ctx)

	var req budgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	budget, err := b.store.CreateBudget(ctx, &model.Budget{
		OwnerID: claims.UserID,
		Name:    req.Name,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create budget")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, budget)
}
