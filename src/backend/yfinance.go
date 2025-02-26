package backend

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Chart struct {
		Result []struct {
			Meta       any     `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close  []float64 `json:"close"`
					Low    []float64 `json:"low"`
					High   []float64 `json:"high"`
					Open   []float64 `json:"open"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
				AdjClose []struct {
					AdjClose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error any
	} `json:"chart"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(symbol string, target interface{}) error {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?1d&interval=1d", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func GetCurrentStockPrice(symbol string) (float32, error) {
	res := &Response{}
	err := getJson(symbol, res)
	if err != nil {
		return -1, err
	}

	// no price returned from backend
	if len(res.Chart.Result) == 0 || len(res.Chart.Result[0].Indicators.AdjClose[0].AdjClose) == 0 {
		return 0, nil
	}
	return float32(res.Chart.Result[0].Indicators.AdjClose[0].AdjClose[0]), nil
}
