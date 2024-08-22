package budgets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (b *Budgets) HandleGetBudgets(ctx *gin.Context) {
	claims := b.auth.GetClaimsFromContext(ctx)
	filter := &store.GetBudgetsFilter{Owner: claims.UserID}

	budgets, err := b.store.GetBudgets(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get budgets")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newBudgetsResponse(budgets))
}
