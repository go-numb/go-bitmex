package realtime_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/go-numb/go-bitmex/realtime"

	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	str := "test:btcusd"
	assert.EqualValues(t, "btcusd", str[strings.Index(str, ":")+1:])
}

func TestParse(t *testing.T) {
	data := `{"ask":[[7024.5,200430],[7024.5,200430]]}`
	var in map[string]interface{}
	if err := json.Unmarshal([]byte(data), &in); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", reflect.TypeOf(in["ask"]))

	body, ok := in["ask"].([]interface{})
	if !ok {
		t.Fatal("type defined")
	}

	books := make([]realtime.Book, len(body))
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
	fmt.Printf("success: %+v\n", books)
}

func TestConnect(t *testing.T) {
	// Only public
	// ctx := context.Background()
	// Both
	ctx := realtime.NewAuth(os.Getenv("MEXKEY"), os.Getenv("MEXSECRET"))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := make(chan realtime.Response, 10)
	public := []string{
		"trade",
	}
	private := []string{
		"order",
	}
	symbols := []string{
		"XBTUSD",
		"ETHUSD",
	}

	go realtime.Connect(ctx, ch, public, symbols, nil)
	go realtime.Connect(ctx, ch, private, nil, nil)

	for {
		select {
		case v := <-ch:
			switch v.Types {
			case realtime.OrderbookL:
				fmt.Printf("%+v\n", v.OrderbookL)

			case realtime.Orderbook:
				fmt.Printf("%+v\n", v.Orderbook)

			case realtime.Quote:
				fmt.Printf("%+v\n", v.Quote)

			case realtime.TradeBin:
				fmt.Printf("%+v\n", v.TradeBin)

			case realtime.Trade:
				fmt.Printf("%+v\n", v.Trade)

			case realtime.Settlement:
				fmt.Printf("%+v\n", v.Settlement)

			case realtime.Funding:
				fmt.Printf("%+v\n", v.Funding)

			case realtime.Chat:
				fmt.Printf("%+v\n", v.Chat)

			case realtime.Announcement:
				fmt.Printf("%+v\n", v.Announcement)

			case realtime.Connected:
				fmt.Printf("%+v\n", v.Connected)

			case realtime.Instrument:
				fmt.Printf("%+v\n", v.Instrument)

			case realtime.Insurance:
				fmt.Printf("%+v\n", v.Insurance)

				/*
					# Private
				*/
			case realtime.Wallet:
				fmt.Printf("%+v\n", v.Wallet)
			case realtime.Margin:
				fmt.Printf("%+v\n", v.Margin)
			case realtime.Execution:
				fmt.Printf("%+v\n", v.Execution)
			case realtime.Order:
				fmt.Printf("%+v\n", v.Order)
			case realtime.Position:
				fmt.Printf("%+v\n", v.Position)

			}

		case <-ctx.Done():
			break
		}
	}
}
