package api

import (
	"github.com/ndovnar/family-budget-api/internal/api/handlers/accounts"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/budgets"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/categories"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/tokens"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/users"
	"github.com/ndovnar/family-budget-api/internal/api/middlewares"
)

func (api *API) registerRoutes() {
	userHandlers := users.New(api.auth, api.store)
	accountHandlers := accounts.New(api.auth, api.store)
	tokenHandlers := tokens.New(api.auth, api.store)
	budgetHandlers := budgets.New(api.auth, api.store)
	categoryHandlers := categories.New(api.store)

	api.router.Use(middlewares.Error())
	authRoutes := api.router.Group("/").Use(middlewares.Auth(api.auth, api.store))

	api.router.POST("/users", userHandlers.HandleCreateUser)
	api.router.POST("/users/login", userHandlers.HandleLoginUser)
	authRoutes.POST("/users/logout", userHandlers.HandleLogoutUser)

	api.router.POST("/tokens/renew/access", tokenHandlers.HandleRenewAccessToken)
	api.router.POST("/tokens/renew/refresh", tokenHandlers.HandleRenewRefreshToken)

	authRoutes.GET("/accounts", accountHandlers.HandleGetAccounts)
	authRoutes.GET("/accounts/:id", accountHandlers.HandleGetAcccount)
	authRoutes.POST("/accounts", accountHandlers.HandleCreateAccount)
	authRoutes.PUT("/accounts/:id", accountHandlers.HandleUpdateAccount)
	authRoutes.DELETE("/accounts/:id", accountHandlers.HandleDeleteAccount)

	authRoutes.GET("/budgets", budgetHandlers.HandleGetBudgets)
	authRoutes.GET("/budgets/:id", budgetHandlers.HandleGetBudget)
	authRoutes.POST("/budgets", budgetHandlers.HandleCreateBudget)
	authRoutes.PUT("/budgets/:id", budgetHandlers.HandleUpdateBudget)
	authRoutes.DELETE("/budgets/:id", budgetHandlers.HandleDeleteBudget)

	authRoutes.GET("/categories", categoryHandlers.HandleGetCategories)
	authRoutes.GET("/categories/:id", categoryHandlers.HandleGetCategory)
	authRoutes.POST("/categories", categoryHandlers.HandleCreateCategory)
	authRoutes.PUT("/categories/:id", categoryHandlers.HandleUpdateCategory)
	authRoutes.DELETE("/categories/:id", categoryHandlers.HandleDeleteCategory)
}
