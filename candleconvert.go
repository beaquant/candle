package candle

import (
	"fmt"
	"github.com/nntaoli-project/GoEx"
	"math"
	"time"
)

const (
	DAY = iota
	HOURS
	MINUTES
)

var (
	isFirstFind = true
	firstStamp  = int64(0)
)

func calcHigh(ks []goex.Kline, n int, baseCycle, newCycle int64) float64 {
	var max = ks[n].High
	for i := 1; i < int(newCycle/baseCycle); i++ {
		max = math.Max(ks[n+i].High, max)
	}
	return max
}

func calcLow(ks []goex.Kline, n int, baseCycle, newCycle int64) float64 {
	var min = ks[n].Low
	for i := 1; i < int(newCycle/baseCycle); i++ {
		min = math.Max(ks[n+i].Low, min)
	}
	return min
}

func getDHM(objTime time.Time, baseCycle, newCycle int64) []int {
	var ret = make([]int, 3)
	if baseCycle%(60*60*24) == 0 {
		ret[0] = objTime.Day()
		ret[1] = DAY
	} else if baseCycle%(60*60) == 0 {
		ret[0] = objTime.Hour()
		ret[1] = HOURS
	} else if baseCycle%(60) == 0 {
		ret[0] = objTime.Minute()
		ret[1] = MINUTES
	}
	if newCycle%(60*60*24) == 0 {
		ret[2] = DAY
	} else if newCycle%(60*60) == 0 {
		ret[2] = HOURS
	} else if newCycle%(60) == 0 {
		ret[2] = MINUTES
	}
	return ret
}

func searchFirstTime(ret []int, baseCycle, newCycle int64) bool {
	if ret[1] == DAY && ret[2] == DAY {
		var array_day []int
		for i := 1; i < 29; i += int(newCycle / baseCycle) {
			array_day = append(array_day, i)
		}
		for j := 0; j < len(array_day); j++ {
			if ret[0] == array_day[j] {
				return true
			}
		}
	} else if ret[1] == HOURS && ret[2] == HOURS {
		var array_hours []int
		for i := 0; i < 24; i += int(newCycle / baseCycle) {
			array_hours = append(array_hours, i)
		}
		for j := 0; j < len(array_hours); j++ {
			if ret[0] == array_hours[j] {
				return true
			}
		}
	} else if ret[1] == MINUTES && ret[2] == MINUTES {
		var array_minutes []int
		for i := 0; i < 60; i += int(newCycle / baseCycle) {
			array_minutes = append(array_minutes, i)
		}
		for j := 0; j < len(array_minutes); j++ {
			if ret[0] == array_minutes[j] {
				return true
			}
		}
	} else {
		panic(fmt.Sprintln("目标周期与基础周期不匹配！目标周期秒数：", newCycle, " 基础周期秒数: ", baseCycle))
	}
	return false
}

func ConvertRecords(records []goex.Kline, newCycle int64) []goex.Kline {
	var AfterAssRecords []goex.Kline

	if len(records) < 2 {
		panic(fmt.Sprintln("传入的records参数为 错误! 基础K线长度小于2"))
	}
	var baseCycle = records[len(records)-1].Timestamp - records[len(records)-2].Timestamp

	if newCycle%baseCycle != 0 {
		panic(fmt.Sprintln("目标周期‘", newCycle, "’不是 基础周期 ‘", baseCycle, "’ 的整倍数，无法合成！"))
	}
	if int(newCycle/baseCycle) > len(records) {
		panic(fmt.Sprintln("基础K线数量不足，请检查是否基础K线周期过小！"))
	}

	// 判断时间戳, 找到 基础K线  相对于 目标K线的起始时间。
	var objTime time.Time
	for i := 0; i < len(records); i++ {
		objTime = time.Unix(records[i].Timestamp, 0)
		var ret = getDHM(objTime, baseCycle, newCycle)

		if isFirstFind == true && searchFirstTime(ret, baseCycle, newCycle) == true {
			firstStamp = records[i].Timestamp
			records = records[i:] // 把目标K线周期前不满足合成的数据排除。
			isFirstFind = false
			break // 排除后跳出
		} else if isFirstFind == false {
			if (records[i].Timestamp-firstStamp)%newCycle == 0 {
				records = records[i:] // 把目标K线周期前不满足合成的数据排除。
				break
			}
		}
	}

	var n = 0
	for n = 0; n < len(records)-int(newCycle/baseCycle); n += int(newCycle / baseCycle) { // 合成
		var BarObj goex.Kline
		BarObj.Pair = records[n].Pair
		BarObj.Timestamp = records[n].Timestamp
		BarObj.Open = records[n].Open
		BarObj.High = calcHigh(records, n, baseCycle, newCycle)
		BarObj.Low = calcLow(records, n, baseCycle, newCycle)
		BarObj.Close = records[n+int(newCycle/baseCycle)-1].Close
		BarObj.Vol = records[n+int(newCycle/baseCycle)-1].Vol
		AfterAssRecords = append(AfterAssRecords, BarObj)
	}
	var BarObj goex.Kline
	fmt.Println("n:", n, "newCycle/baseCycle:", newCycle/baseCycle)
	fmt.Println("baseCycle:", baseCycle, "newCycle:", newCycle)
	BarObj.Timestamp = records[n-int(newCycle/baseCycle)].Timestamp + newCycle // 最后一根时间不能变，
	BarObj.Open = records[n].Open
	BarObj.Close = records[len(records)-1].Close
	BarObj.Vol = records[len(records)-1].Vol
	BarObj.Pair = records[n].Pair
	var max = records[n].High
	var min = records[n].Low
	for index_n := n + 1; index_n < len(records); index_n++ {
		max = math.Max(max, records[index_n].High)
		min = math.Min(min, records[index_n].Low)
	}
	BarObj.High = max
	BarObj.Low = min
	AfterAssRecords = append(AfterAssRecords, BarObj)

	return AfterAssRecords
}
