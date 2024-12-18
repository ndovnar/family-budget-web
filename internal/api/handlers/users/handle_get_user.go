package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (u *Users) HandleGetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	hasAccess := u.authz.IsUserHasAccessToUser(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	user, err := u.store.GetUserByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}

	ctx.JSON(http.StatusOK, user)
}
