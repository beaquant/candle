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

	if candle.Period.Since(ts.LastCandle().Period) != ts.LastCandle().Period.Length() {
		return false
	}
	ts.LastCandle() = candle
	return true
}

func (ts *TimeSeries) PeriodConvert(newPeriod techan.TimePeriod) (t *TimeSeries) {
	if len(ts.Candles)-1 == 0 || ts.LastCandle() == nil {
		panic(fmt.Errorf("Error convert Candle: candle cannot be nil"))
	}

	return
}
