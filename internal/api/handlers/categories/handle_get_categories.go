package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/rs/zerolog/log"
)

func (c *Categories) HandleGetCategories(ctx *gin.Context) {
	filter := &filter.GetCategoriesFilter{}
	if err := ctx.ShouldBindQuery(filter); err != nil {
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	categories, count, err := c.store.GetCategories(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to get categories")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newCategoriesResponse(categories, count))
}
