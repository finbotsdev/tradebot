// +build unit

package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategyBuy2 struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategyBuy2 struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategyBuy2 struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_Buy2(t *testing.T) {
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyBuy2),
		"BollingerBands":                     new(MockBBands_StrategyBuy2),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyBuy2)}
	balances := []common.Coin{
		common.Coin{
			Currency:  "BTC",
			Available: 1.5},
		common.Coin{
			Currency:  "USD",
			Available: 0.0}}
	lastTrade := &common.Trade{
		ID:       1,
		ChartID:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    8000}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Indicators:   strategyIndicators,
		NewPrice:     11000,
		TradeFee:     .025,
		Balances:     balances,
		LastTrade:    lastTrade}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	requiredIndicators := strategy.GetRequiredIndicators()
	assert.Equal(t, "RelativeStrengthIndex", requiredIndicators[0])
	assert.Equal(t, "BollingerBands", requiredIndicators[1])
	assert.Equal(t, "MovingAverageConvergenceDivergence", requiredIndicators[2])

	buy, sell, err := strategy.GetBuySellSignals()
	assert.Equal(t, true, buy)
	assert.Equal(t, false, sell)
	assert.Equal(t, "Out of USD funding!", err.Error())
}

func (mrsi *MockRSI_StrategyBuy2) Calculate(price float64) float64 {
	return 29.0
}

func (mrsi *MockRSI_StrategyBuy2) IsOverBought(rsiValue float64) bool {
	return false
}

func (mrsi *MockRSI_StrategyBuy2) IsOverSold(rsiValue float64) bool {
	return true
}

func (mrsi *MockBBands_StrategyBuy2) Calculate(price float64) (float64, float64, float64) {
	return 14000.0, 13000.0, 12000.0
}