package store

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Store interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error

	GetAccounts(ctx context.Context, filter *filter.GetAccountsFilter) ([]*model.Account, int64, error)
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error)
	DeleteAccount(ctx context.Context, id string) error

	GetBudgets(ctx context.Context, filter *filter.GetBudgetsFilter) ([]*model.Budget, int64, error)
	GetBudget(ctx context.Context, id string) (*model.Budget, error)
	CreateBudget(ctx context.Context, budget *model.Budget) (*model.Budget, error)
	UpdateBudget(ctx context.Context, id string, budget *model.Budget) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id string) error

	GetCategories(ctx context.Context, filter *filter.GetCategoriesFilter) ([]*model.Category, int64, error)
	GetCategory(ctx context.Context, id string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, id string) error

	GetTransactions(ctx context.Context, filter *filter.GetTransactionsFilter) ([]*model.Transaction, int64, error)
	GetTransaction(ctx context.Context, id string) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}
