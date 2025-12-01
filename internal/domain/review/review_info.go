package review

type ReviewInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Text        string `json:"text"`
	ReviewYear  int    `json:"review_year"`
	ReviewMonth int    `json:"review_month"`
	ReviewDay   int    `json:"review_day"`
	UserRating  int    `json:"user_rating"`
	IsLiked     bool   `json:"is_liked"`
	Likes       int    `json:"likes"`
}
