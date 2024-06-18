package mongo

import (
	"context"
	"time"

	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store/mongo/dto"
	"github.com/rs/zerolog/log"
)

func (m *Mongo) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	currentTime := time.Now()

	dtoAccount := dto.ModelAccountToDtoAccount(account)
	dtoAccount.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}
	dtoAccount.IsDeleted = false

	_, err := m.database.
		Collection(CollectionAccounts).
		InsertOne(ctx, dtoAccount)

	if err != nil {
		log.Info().Msgf("mongo: failed to create account. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	createdAccount := dto.DtoAccountToModelAccount(dtoAccount)

	return createdAccount, nil
}
