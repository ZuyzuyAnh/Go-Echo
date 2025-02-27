package dto

type CreateMovieReq struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int64  `json:"duration"`
	CoverURL      string `json:"cover_url"`
	BackgroundURL string `json:"background_url"`
}

type MovieResp struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int64  `json:"duration"`
	CoverURL      string `json:"cover_url"`
	BackgroundURL string `json:"background_url"`
}
