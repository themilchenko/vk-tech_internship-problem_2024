package httpModels

type SortBy string

var EmptyModel = []byte("{}")

type ID struct {
	ID uint64 `json:"id"`
}
