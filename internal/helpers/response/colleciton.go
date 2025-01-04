package response

type CollectionResponse[T any] struct {
	Values []T                 `json:"values"`
	Meta   *CollectionMetaData `json:"_meta"`
}

type CollectionMetaData struct {
	Count int64 `json:"count"`
}
