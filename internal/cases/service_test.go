package cases

type Sevice interface {
	GetCurrencyRate(coinName string) (*models.Coin, error)
	Update
}
