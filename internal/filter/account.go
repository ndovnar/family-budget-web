package filter

type GetAccountsFilter struct {
	OwnerID     string
	Deleted    bool
	Pagination *Pagination
}
