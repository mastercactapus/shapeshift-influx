package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type marketInfo struct {
	Pair      string
	Rate      float64
	Limit     float64
	Min       float64
	MinerFee  float64
	Timestamp time.Time
}

type transaction struct {
	CurIn     string
	CurOut    string
	Amount    float64
	Timestamp float64
}

func floatStr(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.ContainsRune(s, '.') {
		s += ".0"
	}
	return s
}

func (m marketInfo) Entry() string {
	return "market_info,pair=" + m.Pair +
		" rate=" + floatStr(m.Rate) +
		",limit=" + floatStr(m.Limit) +
		",min=" + floatStr(m.Min) +
		",miner_fee=" + floatStr(m.MinerFee) +
		" " + strconv.FormatInt(m.Timestamp.UnixNano(), 10)
}

func (t transaction) Entry() string {
	return "transaction,currency_in=" + t.CurIn + ",currency_out=" + t.CurOut +
		" amount=" + floatStr(t.Amount) +
		" " + strconv.FormatInt(int64(t.Timestamp*1000000000.0), 10)
}

func getPair(pair string) (*marketInfo, error) {
	resp, err := http.Get("https://shapeshift.io/marketinfo/" + pair)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 status code: %s", resp.Status)
	}
	var info marketInfo
	info.Timestamp = time.Now()
	err = json.NewDecoder(resp.Body).Decode(&info)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getTx() ([]transaction, error) {
	resp, err := http.Get("https://shapeshift.io/recenttx/50")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 status code: %s", resp.Status)
	}
	txns := make([]transaction, 0, 50)
	err = json.NewDecoder(resp.Body).Decode(&txns)
	if err != nil {
		return nil, err
	}

	return txns, nil
}
