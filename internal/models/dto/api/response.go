package api

type MessageResponse struct {
	Message string "json:\"message\""
}

type Response[T any] struct {
	Message string "json:\"message\""
	Data    *T     "json:\"data,omitempty\""
}

type Pagination struct {
	Page  int   "json:\"page\""
	Limit int   "json:\"limit\""
	Total int64 "json:\"total\""
}

type PaginatedResponse[T any] struct {
	Data       []T        "json:\"data\""
	Pagination Pagination "json:\"pagination\""
}
