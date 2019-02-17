package candle

import (
	"fmt"
	"github.com/sdcoffey/techan"
)

type TimeSeries struct {
	techan.TimeSeries
}

// ModifyLastCandle modify this latest candle of TimeSeries with the given candle.
func (ts *TimeSeries) ModifyLastCandle(candle *techan.Candle) bool {
	if candle == nil || ts.LastCandle() == nil {
		panic(fmt.Errorf("Error modify Candle: candle cannot be nil"))
	}
	if candle.Period.Since(ts.LastCandle().Period) != ts.LastCandle().Period.Length() {
		return false
	}
	ts.LastCandle() = candle
	return true
}

// UpdateLastCandle update this latest candle of TimeSeries with the given candle.
func (ts *TimeSeries) UpdateLastCandle(candle *techan.Candle) bool {
	if candle == nil || ts.LastCandle() == nil {
		panic(fmt.Errorf("Error modify Candle: candle cannot be nil"))
	}

	if (ts.LastCandle().Period.Start.Before(candle.Period.Start) || ts.LastCandle().Period.Start.Equal(candle.Period.Start)) &&
		(ts.LastCandle().Period.End.After(candle.Period.End) || ts.LastCandle().Period.End.Equal(candle.Period.End)) {
		ts.LastCandle().Volume = ts.LastCandle().Volume.Add(candle.Volume)
		ts.LastCandle().TradeCount += candle.TradeCount
		ts.LastCandle().ClosePrice = candle.ClosePrice

		if ts.LastCandle().OpenPrice.Zero() {
			ts.LastCandle().OpenPrice = candle.OpenPrice
		}

		if ts.LastCandle().MaxPrice.Zero() {
			ts.LastCandle().MaxPrice = candle.MaxPrice
		} else if candle.MaxPrice.GT(ts.LastCandle().MaxPrice) {
			ts.LastCandle().MaxPrice = candle.MaxPrice
		}

		if ts.LastCandle().MinPrice.Zero() {
			ts.LastCandle().MinPrice = candle.MinPrice
		} else if candle.MinPrice.LT(ts.LastCandle().MinPrice) {
			ts.LastCandle().MinPrice = candle.MinPrice
		}

		if ts.LastCandle().Volume.Zero() {
			ts.LastCandle().Volume = candle.Volume
		} else {
			ts.LastCandle().Volume = ts.LastCandle().Volume.Add(candle.Volume)
		}
	} else {
		return ts.AddCandle(candle)
	}
	return true
}

func (ts *TimeSeries) PeriodConvert(newPeriod techan.TimePeriod) (t *TimeSeries) {
	if len(ts.Candles)-1 == 0 || ts.LastCandle() == nil {
		panic(fmt.Errorf("Error convert Candle: candle cannot be nil"))
	}

	return
}
