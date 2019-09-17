package bitmex

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	// APIREMAIN is API Limit initial number
	APIREMAIN = 60 // par 60sec
)

// Limit is API Limit struct
type Limit struct {
	Limit  int       // Limit is resets count
	Remain int       // Remain is 残Requests
	Reset  time.Time // Reset Remainの詳細時間(sec未満なし)
}

// NewLimit is API Limit
func NewLimit(isPrivate bool) *Limit {
	if isPrivate {
		return &Limit{
			Limit:  APIREMAIN,
			Remain: APIREMAIN,
			Reset:  time.Now().Add(time.Minute),
		}
	}

	return &Limit{
		Limit:  APIREMAIN,
		Remain: APIREMAIN,
		Reset:  time.Now().Add(time.Minute),
	}
}

// FromHeader X-xxxからLimitを取得
func (p *Limit) FromHeader(h http.Header) {
	period := h.Get("x-ratelimit-limit") // リセット後の残回数
	if period != "" {
		p.Limit, _ = strconv.Atoi(period)
	}
	remain := h.Get("x-ratelimit-remaining") // 残回数
	if remain != "" {
		p.Remain, _ = strconv.Atoi(remain)
	}
	t := h.Get("X-ratelimit-reset") // リセットUTC時間(sec未満なし)
	if t != "" {
		reset, _ := strconv.ParseInt(t, 10, 64)
		p.toTime(reset)
	}
}

// Check is checks remain number
func (p *Limit) Check() error {
	if p.Remain <= 0 {
		if time.Now().After(p.Reset) { // APIRESET時間を過ぎていたらRemainを補充
			p.Remain = p.Limit
		}
		return fmt.Errorf("api limit, has API Limit Remain:%d, Reset time: %s(%s)",
			p.Remain,
			p.Reset.Format("15:04:05"),
			time.Now().Format("15:04:05"))
	}
	return nil
}

// int64 to time.Time
func (p *Limit) toTime(t int64) {
	p.Reset = time.Unix(t, 10)
}
