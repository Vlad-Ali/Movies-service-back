package reviewservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/review/modelconfig"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	"github.com/ollama/ollama/api"
)

type ReviewProvider struct {
	reviewService reviewdomain.Service
	client        *api.Client
	config        modelconfig.ModelConfig
}

func NewReviewProvider(reviewService reviewdomain.Service, config modelconfig.ModelConfig) *ReviewProvider {
	baseURL, err := url.Parse(config.OllamaHost)
	if err != nil {
		slog.Error("Error while parsing the URL ", "Error", err)
		os.Exit(1)
	}

	client := api.NewClient(baseURL, http.DefaultClient)

	return &ReviewProvider{reviewService: reviewService, config: config, client: client}
}

func (r *ReviewProvider) ProvideMovieReviews(ctx context.Context, movieInfo object2.MovieInfo) (string, error) {
	reviews, err := r.reviewService.GetUserReviews(ctx, movieInfo)

	if err != nil {
		return "", err
	}

	if len(reviews) == 0 {
		return "", nil
	}

	var texts []string
	for _, review := range reviews {
		texts = append(texts, review.Text)
	}

	reviewsText := strings.Join(texts, "\n")

	finalUserPrompt := fmt.Sprintf(r.config.UserPrompt, movieInfo.Title, reviewsText)

	request := &api.ChatRequest{
		Model: r.config.Name,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: r.config.SystemPrompt,
			},
			{
				Role:    "user",
				Content: finalUserPrompt,
			},
		},
		Stream: boolPtr(false),
		Options: map[string]interface{}{
			"num_predict": 100,
			"temperature": 0.1,
		},
	}

	var response string
	slog.Debug("Got reviews", "prompt", finalUserPrompt)
	newCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	err = r.client.Chat(newCtx, request, func(resp api.ChatResponse) error {
		response = resp.Message.Content
		return nil
	})

	slog.Debug("Got response", "prompt", response)
	if err != nil {
		slog.Error("Error while creating chat completion", "Error", err)
		return "", err
	}

	if response == "" {
		slog.Warn("Empty response from Ollama")
		return "", nil
	}

	if response[len(response)-1] != '.' {
		splitResponse := strings.Split(response, ".")
		response = strings.Join(splitResponse[:len(splitResponse)-1], ".")
	}

	return response, nil
}

func boolPtr(b bool) *bool {
	return &b
}
