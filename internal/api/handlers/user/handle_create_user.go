package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/hash"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (user *User) HandleCreateUser(ctx *gin.Context) {
	log.Debug().Msg("create user: handle")

	var req newUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("create user: failed to parse data")
		ctx.Error(error.NewHttpError(http.StatusBadRequest))
		return
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("create user: failed to hash password")
		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	newUser, err := user.store.CreateUser(ctx, &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashedPassword,
	})
	if err != nil {
		log.Error().Err(err).Msg("create user: failed to create")

		if err == store.ErrDuplicateKey {
			ctx.Error(error.NewHttpError(http.StatusForbidden))
			return
		}

		ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		return
	}

	log.Debug().Msg("create user: success")
	ctx.JSON(http.StatusOK, newUserResponse(newUser))
}
