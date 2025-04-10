package application

import (
	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidRecordsPrices(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"100.5"}}}},
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"200.75"}}}},
	}

	expected := []float64{100.5, 200.75}
	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestEmptyRecordsReturnsEmpty(t *testing.T) {
	records := []domain.OrderBookRecord{}

	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestInvalidPriceFormatError(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"invalid"}}}},
	}

	result, err := extractPrices(records)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoBidUpdatesSkipsRecord(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{}}},
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"150.25"}}}},
	}

	expected := []float64{150.25}
	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestShortSMAGreaterThanLongSMABuy(t *testing.T) {
	lastShortSMA := 105.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Buy, result)
}

func TestShortSMALessThanLongSMASell(t *testing.T) {
	lastShortSMA := 95.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Sell, result)
}

func TestShortSMAEqualsLongSMANeutral(t *testing.T) {
	lastShortSMA := 100.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Neutral, result)
}

func TestSMAValidInputCorrectSMA(t *testing.T) {
	prices := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	period := 3

	expected := []float64{2, 3, 4, 5, 6, 7, 8, 9}
	result := calculateSMA(prices, period)

	assert.Equal(t, expected, result)
}

func TestSMAPeriodGreaterThanPricesNil(t *testing.T) {
	prices := []float64{1, 2, 3}
	period := 5

	result := calculateSMA(prices, period)

	assert.Nil(t, result)
}

func TestSMAEmptyPricesNil(t *testing.T) {
	prices := []float64{}
	period := 3

	result := calculateSMA(prices, period)

	assert.Nil(t, result)
}

func TestSMAPeriodEqualsPricesLengthSingleValue(t *testing.T) {
	prices := []float64{1, 2, 3, 4, 5}
	period := 5

	expected := []float64{3}
	result := calculateSMA(prices, period)

	assert.Equal(t, expected, result)
}

func SelectSignalWhenShortSMAIsGreaterThanLongSMAReturnsBuy(t *testing.T) {
	lastShortSMA := 105.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Buy, result)
}

func SelectSignalWhenShortSMAIsLessThanLongSMAReturnsSell(t *testing.T) {
	lastShortSMA := 95.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Sell, result)
}

func SelectSignalWhenShortSMAEqualsLongSMAReturnsNeutral(t *testing.T) {
	lastShortSMA := 100.0
	lastLongSMA := 100.0

	result := selectSignal(lastShortSMA, lastLongSMA)

	assert.Equal(t, domain.Neutral, result)
}

func ExtractPricesWithValidRecordsReturnsPrices(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"100.5"}}}},
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"200.75"}}}},
	}

	expected := []float64{100.5, 200.75}
	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func ExtractPricesWithEmptyRecordsReturnsEmptySlice(t *testing.T) {
	records := []domain.OrderBookRecord{}

	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func ExtractPricesWithInvalidPriceFormatReturnsError(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"invalid"}}}},
	}

	result, err := extractPrices(records)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func ExtractPricesWithNoBidUpdatesSkipsRecord(t *testing.T) {
	records := []domain.OrderBookRecord{
		{Data: domain.OrderBookData{BidUpdates: [][]string{}}},
		{Data: domain.OrderBookData{BidUpdates: [][]string{{"150.25"}}}},
	}

	expected := []float64{150.25}
	result, err := extractPrices(records)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
