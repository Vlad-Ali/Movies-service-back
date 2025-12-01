package response

type GetReviewResponse struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	ReviewYear  int    `json:"review_year"`
	ReviewMonth int    `json:"review_month"`
	ReviewDay   int    `json:"review_day"`
}
