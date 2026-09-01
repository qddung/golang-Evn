package bookmark_model

type NewBookmarkRequest struct {
	Url         string `json:"url" validate:"url,required"`
	Description string `json:"description"`
}
