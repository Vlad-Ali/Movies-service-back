package response

type GetReviewResponse struct {
	Text        string `json:"text"`
	ReviewYear  int    `json:"review_year"`
	ReviewMonth int    `json:"review_month"`
	ReviewDay   int    `json:"review_day"`
}
