package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Client struct {
	apiKey string
	Client *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "api key is required")
	}

	return &Client{
		apiKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Client) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	url := fmt.Sprintf("http://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USD", joinTitles(titles))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Apikey %s", c.apiKey))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(entities.ErrInvalidParam, "unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	if len(result) == 0 {
		return nil, errors.Wrap(entities.MissingData, "empty response from API")
	}

	var coins []entities.Coin
	for title, rates := range result {
		coins = append(coins, entities.Coin{
			CoinName: title,
			Price:    rates["USD"],
		})
	}

	return coins, nil
}

func joinTitles(titles []string) string {
	result := ""
	for i, title := range titles {
		if i > 0 {
			result += ","
		}
		result += title
	}
	return result
}
