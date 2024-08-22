package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/rs/zerolog/log"
)

func (a *Categories) HandleGetCategories(ctx *gin.Context) {
	categories, err := a.store.GetCategories(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get categories")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newCategoriesResponse(categories))
}
