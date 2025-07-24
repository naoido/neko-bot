package zenn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiBaseURL = "https://zenn.dev/api"

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"published_at"`
	Path        string    `json:"path"`
	User        User      `json:"user"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarSmallURL string `json:"avatar_small_url"`
}

type articlesResponse struct {
	Articles []Article `json:"articles"`
}

type Repository interface {
	FetchArticlesByUsername(ctx context.Context, username string) ([]Article, error)
}

type repository struct {
	client *http.Client
}

func NewRepository() Repository {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return &repository{
		client: client,
	}
}

func (r *repository) FetchArticlesByUsername(ctx context.Context, username string) ([]Article, error) {
	url := fmt.Sprintf("%s/articles?username=%s&order=latest", apiBaseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res articlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return res.Articles, nil
}
