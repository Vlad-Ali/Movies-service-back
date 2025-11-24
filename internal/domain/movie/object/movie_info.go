package object

import (
	"net/http"
	"strconv"
)

type MovieInfo struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
}

func GetMovieInfoFromReq(r *http.Request) (MovieInfo, error) {
	var movieInfo MovieInfo
	movieInfo.Title = r.URL.Query().Get("title")
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		return movieInfo, err
	}

	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		return movieInfo, err
	}

	day, err := strconv.Atoi(r.URL.Query().Get("day"))
	if err != nil {
		return movieInfo, err
	}

	movieInfo.Year = year
	movieInfo.Month = month
	movieInfo.Day = day
	return movieInfo, nil
}
