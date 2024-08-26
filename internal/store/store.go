package store

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/model"
)

type Store interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error

	GetAccounts(ctx context.Context, filter *model.GetAccountsFilter) ([]*model.Account, error)
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error)
	DeleteAccount(ctx context.Context, id string) error

	GetBudgets(ctx context.Context, filter *model.GetBudgetsFilter) ([]*model.Budget, error)
	GetBudget(ctx context.Context, id string) (*model.Budget, error)
	CreateBudget(ctx context.Context, budget *model.Budget) (*model.Budget, error)
	UpdateBudget(ctx context.Context, id string, budget *model.Budget) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id string) error

	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetCategory(ctx context.Context, id string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, id string) error

	GetTransactions(ctx context.Context, filter *model.GetTransactionsFilter) ([]*model.Transaction, error)
	GetTransaction(ctx context.Context, id string) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}
