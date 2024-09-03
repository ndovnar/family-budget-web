package authz

import (
	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/store"
)

type Authz struct {
	auth  *auth.Auth
	store store.Store
}

func New(auth *auth.Auth, store store.Store) *Authz {
	return &Authz{
		auth:  auth,
		store: store,
	}
}
