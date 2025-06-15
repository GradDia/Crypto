package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
	logger     *slog.Logger
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

func NewClient(apiKey string, logger *slog.Logger, opts ...ClientOption) (*Client, error) {
	const op = "cryptocompare.NewClient"
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("op", op))

	logger.Debug("Initializing client",
		slog.Bool("has_api_key", apiKey != ""),
		slog.Int("options_count", len(opts)))

	if apiKey == "" {
		err := errors.Wrap(entities.ErrInvalidParam, "api key is required")
		logger.Error("Validation failed",
			slog.String("error", err.Error()))
		return nil, err
	}

	client := &Client{
		apiKey:     apiKey,
		HttpClient: &http.Client{},
		priceIn:    defaultPriceIn,
		logger:     logger,
	}

	client.SetOptions(opts...)

	logger.Info("Initializing new client success",
		slog.String("default_currency", defaultPriceIn))
	return client, nil
}

func (c *Client) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	const op = "cryptocompare.GetActualRates"
	logger := c.logger.With(
		slog.String("op", op),
		slog.Int("titles_count", len(titles)))
	startTime := time.Now()

	logger.Debug("Building API request")
	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
		c.logger.Error("Invalid request", slog.String("error", err.Error()))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", baseUrl, multivalues), nil)
	if err != nil {
		logger.Error("Request creation failed",
			slog.String("error", err.Error()),
			slog.String("url", baseUrl+multivalues))
		return nil, errors.Wrapf(entities.ErrInternal, "new request error: %v", err)
	}

	q := req.URL.Query()
	q.Add(queryTsyms, c.priceIn)
	q.Add(queryFsyms, strings.Join(titles, ","))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Apikey %s", c.apiKey))

	logger.Debug("Sending request",
		slog.String("symbols", strings.Join(titles, ",")),
		slog.String("target_currency", c.priceIn))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		logger.Error("Request failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrapf(entities.ErrInternal, "execute request failure: %v", err)
	}
	defer resp.Body.Close()

	logger.Debug("Response received",
		slog.Int("status_code", resp.StatusCode),
		slog.String("status", resp.Status))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error("API returned error",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response", string(body)))
		return nil, errors.Wrapf(entities.ErrInvalidParam, "unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]map[string]float64
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Response parsing failed",
			slog.String("error", err.Error()))
		return nil, errors.Wrapf(entities.ErrInternal, "readAll resp Body: %v", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "unmarshal body error :%v", err)
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

	logger.Info("Request completed successfully",
		slog.Int("coins_received", len(coins)),
		slog.Duration("duration", time.Since(startTime)))

	return coins, nil
}
