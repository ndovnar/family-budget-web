package api

import (
	"github.com/ndovnar/family-budget-api/internal/api/handlers/account"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/user"
	"github.com/ndovnar/family-budget-api/internal/api/middlewares"
)

func (api *API) registerRoutes() {
	userHandlers := user.NewUser(api.auth, api.store)
	accountHandlers := account.NewAccount(api.store)

	api.router.Use(middlewares.Error())
	api.router.POST("/users", userHandlers.HandleCreateUser)
	api.router.POST("/login", userHandlers.HandleLoginUser)

	authRoutes := api.router.Group("/").Use(middlewares.Auth(api.auth, api.store))
	authRoutes.POST("/account", accountHandlers.HandleCreateAccount)
}
