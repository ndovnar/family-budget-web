package filter

type GetBudgetsFilter struct {
	OwnerID    string
	Deleted    bool
	Pagination *Pagination
}
