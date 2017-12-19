package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Moving%20Averages.htm
func TestExponentialMovingAverage(t *testing.T) {
	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(64.75)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.79)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.73)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.73)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.55)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.19)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.91)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.85)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(62.95)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(63.37)})

	ema := NewExponentialMovingAverage(candlesticks)

	/*
		actual := ema.GetMultiplier()
		expected := 0.181818
		if actual != expected {
			t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
		}*/

	actual := ema.GetAverage()
	expected := decimal.NewFromFloat(float64(63.682))
	if !actual.Equals(expected) {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual.Float64(), expected.Float64())
	}

	ema.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(61.33))})
	actual = ema.GetAverage()
	expected = decimal.NewFromFloat(float64(63.254))
	if !actual.Equals(expected) {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(61.51))})
	actual = ema.GetAverage()
	expected = decimal.NewFromFloat(float64(62.937))
	if !actual.Equals(expected) {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(61.87))})
	actual = ema.GetAverage()
	expected = decimal.NewFromFloat(float64(62.743))
	if !actual.Equals(expected) {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

}
