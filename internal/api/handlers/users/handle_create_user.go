package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (u *Users) HandleCreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	user, err := u.store.CreateUser(ctx, &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashedPassword,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")

		if err == store.ErrDuplicateKey {
			ctx.Error(error.NewHttpError(http.StatusForbidden))
			return
		}

		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}
