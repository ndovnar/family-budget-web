package budgets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/rs/zerolog/log"
)

type GetBudgetsQueryParams struct {
	Pagination filter.Pagination
}

func (b *Budgets) HandleGetBudgets(ctx *gin.Context) {
	filter := &filter.GetBudgetsFilter{}
	if err := ctx.ShouldBindQuery(filter); err != nil {
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	claims := b.auth.GetClaimsFromContext(ctx)
	filter.OwnerID = claims.UserID

	budgets, count, err := b.store.GetBudgets(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get budgets")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newBudgetsResponse(budgets, count))
}
