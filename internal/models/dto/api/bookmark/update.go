package bookmark_model

type UpdateBookmarkRequest struct {
	Url         string `json:"url" validate:"url;required"`
	Description string `json:"description"`
}
