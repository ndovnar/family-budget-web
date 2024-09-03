package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (a *Categories) HandleDeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	hasAccess := a.authz.IsUserHasAccessToCategory(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	err := a.store.DeleteCategory(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete category")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}
}
