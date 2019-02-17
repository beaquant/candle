package candle

import (
	"fmt"
	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/binance"
	"net/http"
	"testing"
	"time"

	"net"
	"net/url"
)

var httpProxyClient = &http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return &url.URL{
				Scheme: "socks5",
				Host:   "127.0.0.1:1080"}, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
}

func TestAssembleRecords(t *testing.T) {
	exchange := binance.New(httpProxyClient, "", "")
	records, err := exchange.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 60, int(time.Now().Add(-time.Hour).UnixNano()))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range records {
		fmt.Println(time.Unix(v.Timestamp, 0), v)
	}
	newRecords := ConvertRecords(records, 60*30)
	for _, vv := range newRecords {
		fmt.Println(time.Unix(vv.Timestamp, 0), vv)
	}
	fmt.Println("len records:", len(records))
	fmt.Println("len newRecords:", len(newRecords))
}
