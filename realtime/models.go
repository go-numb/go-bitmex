package realtime

import (
	"fmt"
	"time"
)

type OrderBook struct {
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
	Bids      []Book    `json:"bids"`
	Asks      []Book    `json:"asks"`
}

type Book struct {
	Price float64
	Size  float64
}

type In struct {
}

var bookname = []string{"asks", "bids"}

func (p *OrderBook) UnmarshalJSON(b []byte) error {
	var in map[string]interface{}
	if err := json.Unmarshal(b, &in); err != nil {
		return err
	}

	p.Symbol = fmt.Sprintf("%v", in["symbol"])
	p.Timestamp, _ = time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", in["timestamp"]))

	for i := range bookname {
		body, ok := in[bookname[i]].([]interface{})
		if !ok {
			return fmt.Errorf("data has not type")
		}

		books := make([]Book, len(body))
		for i, v := range body {
			book, ok := v.([]interface{})
			if !ok {
				continue
			}

			if len(book) != 2 {
				continue
			}
			books[i].Price, _ = book[0].(float64)
			books[i].Size, _ = book[1].(float64)
		}

		if bookname[i] == "ask" {
			p.Asks = books
		} else {
			p.Bids = books
		}
	}

	return nil
}
