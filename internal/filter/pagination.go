package filter

type Pagination struct {
	Limit  int64 `form:"limit,default=50"`
	Offset int64 `form:"offset,default=0"`
}
