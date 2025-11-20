package object

type MovieInfo struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
}
