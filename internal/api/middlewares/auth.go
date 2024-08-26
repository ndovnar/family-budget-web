package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func Auth(auth *auth.Auth, store store.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			ctx.Error(error.NewHttpErrorWithDescription(http.StatusUnauthorized, "authorization header is not provided"))
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.Error(error.NewHttpErrorWithDescription(http.StatusUnauthorized, "invalid authorization header format"))
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			ctx.Error(error.NewHttpErrorWithDescription(http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType)))
			ctx.Abort()
			return
		}

		accessToken := fields[1]
		claims, err := auth.VerifyToken(accessToken)
		if err != nil {
			log.Error().Err(err).Msg("auth middleware: failed to verify token")
			ctx.Error(error.NewHttpError(http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		session, err := store.GetSessionByID(ctx, claims.SessionID)
		if err != nil {
			log.Error().Err(err).Msg("auth middleware: failed to get session")
			ctx.Error(error.NewHttpError(http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		if session.IsDeleted {
			log.Error().Err(err).Msg("auth middleware: session is revoked")
			ctx.Error(error.NewHttpError(http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		auth.SetClaimsToContext(ctx, claims)
		ctx.Next()
	}
}
