package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"Cryptoproject/internal/entities"
)

const (
	baseUrl     = "https://min-api.cryptocompare.com"
	multivalues = "/data/pricemulti"

	defaultPriceIn = "USD"
	queryFsyms     = "fsyms"
	queryTsyms     = "tsyms"
)

type Client struct {
	apiKey     string
	HttpClient *http.Client
	priceIn    string
}

type ClientOption func(client *Client)

func WithPriceIn(priceIn string) ClientOption {
	return func(c *Client) {
		c.priceIn = priceIn
	}
}

func (c *Client) SetOptions(opts ...ClientOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "api key is required")
	}

	client := &Client{
		apiKey:     apiKey,
		HttpClient: &http.Client{},
		priceIn:    defaultPriceIn,
	}

	client.SetOptions(opts...)

	return client, nil
}

func (c *Client) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", baseUrl, multivalues), nil)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "new request error: %v", err)
	}

	q := req.URL.Query()
	q.Add(queryTsyms, c.priceIn)
	q.Add(queryFsyms, strings.Join(titles, ","))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Apikey %s", c.apiKey))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "execute request failure: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(entities.ErrInvalidParam, "unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]map[string]float64
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	if len(result) == 0 {
		return nil, errors.Wrap(entities.ErrNotFound, "empty response from API")
	}

	var coins []entities.Coin
	for title, rates := range result {
		coins = append(coins, entities.Coin{
			CoinName: title,
			Price:    rates[c.priceIn],
		})
	}

	return coins, nil
}

//Переписать тест
