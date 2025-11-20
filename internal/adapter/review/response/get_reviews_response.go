package response

import "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"

type GetReviewsResponse struct {
	Reviews []*review.ReviewInfo `json:"reviews"`
}
