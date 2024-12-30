package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/accounts"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/budgets"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/categories"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/tokens"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/transactions"
	"github.com/ndovnar/family-budget-api/internal/api/handlers/users"
	"github.com/ndovnar/family-budget-api/internal/api/middlewares"
)

type routeGroup struct {
	path        string
	middlewares []gin.HandlerFunc
	routes      []route
}

type route struct {
	path        string
	method      string
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
}

func (api *API) registerRoutes() {
	tokenHandlers := tokens.New(api.auth, api.store)
	userHandlers := users.New(api.auth, api.authz, api.store)
	accountHandlers := accounts.New(api.auth, api.authz, api.store)
	budgetHandlers := budgets.New(api.auth, api.authz, api.store)
	categoryHandlers := categories.New(api.authz, api.store)
	transactionHandlers := transactions.New(api.auth, api.authz, api.store)

	authMiddleware := middlewares.Auth(api.auth, api.store)
	api.router.Use(middlewares.Error())

	routes := []routeGroup{
		{
			path: "tokens",
			routes: []route{
				{
					path:    "/renew/access",
					method:  http.MethodPost,
					handler: tokenHandlers.HandleRenewAccessToken,
				},
				{
					path:    "/renew/refresh",
					method:  http.MethodPost,
					handler: tokenHandlers.HandleRenewRefreshToken,
				},
			},
		},
		{
			path: "/users",
			routes: []route{
				{
					path:    "",
					method:  http.MethodPost,
					handler: userHandlers.HandleCreateUser,
				},
				{
					path:    "/:id",
					method:  http.MethodGet,
					handler: userHandlers.HandleGetUser,
					middlewares: []gin.HandlerFunc{
						authMiddleware,
					},
				},
				{
					path:    "/login",
					method:  http.MethodPost,
					handler: userHandlers.HandleLoginUser,
				},
				{
					path:    "/logout",
					method:  http.MethodPost,
					handler: userHandlers.HandleLogoutUser,
					middlewares: []gin.HandlerFunc{
						authMiddleware,
					},
				},
			},
		},
		{
			path: "/accounts",
			middlewares: []gin.HandlerFunc{
				authMiddleware,
			},
			routes: []route{
				{
					path:    "",
					method:  http.MethodGet,
					handler: accountHandlers.HandleGetAccounts,
				},
				{
					path:    "/:id",
					method:  http.MethodGet,
					handler: accountHandlers.HandleGetAcccount,
				},
				{
					path:    "",
					method:  http.MethodPost,
					handler: accountHandlers.HandleCreateAccount,
				},
				{
					path:    "/:id",
					method:  http.MethodPut,
					handler: accountHandlers.HandleUpdateAccount,
				},
				{
					path:    "/:id",
					method:  http.MethodDelete,
					handler: accountHandlers.HandleDeleteAccount,
				},
			},
		},
		{
			path: "/budgets",
			middlewares: []gin.HandlerFunc{
				authMiddleware,
			},
			routes: []route{
				{
					path:    "",
					method:  http.MethodGet,
					handler: budgetHandlers.HandleGetBudgets,
				},
				{
					path:    "/:id",
					method:  http.MethodGet,
					handler: budgetHandlers.HandleGetBudget,
				},
				{
					path:    "",
					method:  http.MethodPost,
					handler: budgetHandlers.HandleCreateBudget,
				},
				{
					path:    "/:id",
					method:  http.MethodPut,
					handler: budgetHandlers.HandleUpdateBudget,
				},
				{
					path:    "/:id",
					method:  http.MethodDelete,
					handler: budgetHandlers.HandleDeleteBudget,
				},
			},
		},
		{
			path: "/categories",
			middlewares: []gin.HandlerFunc{
				authMiddleware,
			},
			routes: []route{
				{
					path:    "",
					method:  http.MethodGet,
					handler: categoryHandlers.HandleGetCategories,
				},
				{
					path:    "/:id",
					method:  http.MethodGet,
					handler: categoryHandlers.HandleGetCategory,
				},
				{
					path:    "",
					method:  http.MethodPost,
					handler: categoryHandlers.HandleCreateCategory,
				},
				{
					path:    "/:id",
					method:  http.MethodPut,
					handler: categoryHandlers.HandleUpdateCategory,
				},
				{
					path:    "/:id",
					method:  http.MethodDelete,
					handler: categoryHandlers.HandleDeleteCategory,
				},
			},
		},
		{
			path: "/transactions",
			middlewares: []gin.HandlerFunc{
				authMiddleware,
			},
			routes: []route{
				{
					path:    "",
					method:  http.MethodGet,
					handler: transactionHandlers.HandleGetTransactions,
				},
				{
					path:    "/:id",
					method:  http.MethodGet,
					handler: transactionHandlers.HandleGetTransaction,
				},
				{
					path:    "",
					method:  http.MethodPost,
					handler: transactionHandlers.HandleCreateTransaction,
				},
				{
					path:    "/:id",
					method:  http.MethodPut,
					handler: transactionHandlers.HandleUpdateTransaction,
				},
				{
					path:    "/:id",
					method:  http.MethodDelete,
					handler: transactionHandlers.HandleDeleteTransaction,
				},
			},
		},
	}

	setupRouter(api.router, routes)
}

func setupRouter(router *gin.Engine, routeGroups []routeGroup) {
	for _, group := range routeGroups {
		setupRouteGroup(router, group)
	}
}

func setupRouteGroup(router *gin.Engine, group routeGroup) {
	routerGroup := router.Group(group.path)

	setupMiddlewares(routerGroup, group.middlewares)

	for _, route := range group.routes {
		setupRoute(routerGroup, route)
	}
}

func setupRoute(router *gin.RouterGroup, route route) {
	routerGroup := router.Group(route.path)

	setupMiddlewares(routerGroup, route.middlewares)

	switch route.method {
	case http.MethodGet:
		routerGroup.GET("", route.handler)
	case http.MethodPost:
		routerGroup.POST("", route.handler)
	case http.MethodPut:
		routerGroup.PUT("", route.handler)
	case http.MethodPatch:
		routerGroup.PATCH("", route.handler)
	case http.MethodDelete:
		routerGroup.DELETE("", route.handler)
	}
}

func setupMiddlewares(routerGroup *gin.RouterGroup, middlewares []gin.HandlerFunc) {
	for _, middleware := range middlewares {
		routerGroup.Use(middleware)
	}
}
