package realtime

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-numb/go-bitmex"

	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"
)

const (
	ENDPOINT                   = "wss://www.bitmex.com/realtime"
	READDEADLINE time.Duration = 300 * time.Second

	AUTHKEY = "auth"
)

type Types int

const (
	All          Types = iota
	Announcement       // お知らせ
	Chat
	Connected  // 接続ユーザ/Bot統計
	Funding    // 資金調達率
	Instrument // 商品情報
	Insurance  // 保険基金情報
	Liquidation
	Orderbook
	OrderbookL
	Notifications // システム通知
	Quote         // ブック上位レベル
	Settlement    // 決済
	Trade         // 約定
	TradeBin

	// Private
	Affiliate
	Execution
	Order
	Margin
	Position
	NotificationsForPrivate
	Transact // 入出金に関する情報
	Wallet   // アドレス残高

	Undefined
	Error
)

func (p Types) String() string {
	switch p {
	case Announcement: // お知らせ
		return "announcement"
	case Chat:
		return "chat"
	case Connected: // 接続ユーザ/Bot統計
		return "connected"
	case Funding: // 資金調達率
		return "funding"
	case Instrument: // 商品情報
		return "instrument"
	case Insurance: // 保険基金情報
		return "insurance"
	case Liquidation:
		return "liquidation"
	case Orderbook:
		return "orderBook"
	case OrderbookL:
		return "orderBookL"
	case Notifications: // システム通知
		return "publicNotifications"
	case Quote: // ブック上位レベル
		return "quote"
	case Settlement: // 決済
		return "settlement"
	case Trade: // 約定
		return "trade"
	case TradeBin:
		return "tradeBin"

	// Private
	case Affiliate:
		return "affiliate"
	case Execution:
		return "execution"
	case Order:
		return "order"
	case Margin:
		return "margin"
	case Position:
		return "position"
	case NotificationsForPrivate:
		return "privateNotifications"
	case Transact: // 入出金に関する情報
		return "transact"
	case Wallet: // アドレス残高
		return "wallet"

	case Error:
		return "error"
	}
	return "undefined"
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	conn *websocket.Conn
	Auth *Auth

	log *log.Logger
}

type Auth struct {
	Key    string
	Secret string
}

func New(ctx context.Context, l *log.Logger) *Client {
	if l == nil {
		l = log.New(os.Stdout, "websocket jsonrpc", log.Lmicroseconds)
	}

	auth, ok := ctx.Value(AUTHKEY).(*Auth)
	if !ok {
		auth = nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(ENDPOINT, nil)
	if err != nil {
		return nil
	}
	return &Client{
		conn: conn,
		Auth: auth,
		log:  l,
	}
}

func NewAuth(key, secret string) context.Context {
	return context.WithValue(context.TODO(), AUTHKEY, &Auth{
		Key:    key,
		Secret: secret,
	})
}

func (p *Client) Close() error {
	if p.conn == nil {
		return fmt.Errorf("connection no longer exists")
	}

	if err := p.conn.Close(); err != nil {
		p.log.Println(err)
		return err
	}

	return nil
}

type Request struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
	ID   int           `json:"id,omitempty"`
}

type Response struct {
	Types       Types
	ProductCode string
	Action      string

	Announcement  []bitmex.Announcement
	Chat          []bitmex.Chat
	Connected     []bitmex.ConnectedUsers
	Funding       []bitmex.Funding
	Instrument    []bitmex.Instrument
	Insurance     []bitmex.Insurance
	Settlement    []bitmex.Settlement
	Notifications []bitmex.Notification

	Orderbook  []OrderBook
	OrderbookL []bitmex.OrderBookL2
	Quote      []bitmex.Quote
	Trade      []bitmex.Trade
	TradeBin   []bitmex.TradeBin

	// Private
	Affiliate               []bitmex.Affiliate
	Execution               []bitmex.Execution
	Order                   []bitmex.Order
	Margin                  []bitmex.Margin
	Position                []bitmex.Position
	NotificationsForPrivate []bitmex.Notification
	Transact                []bitmex.Transaction
	Wallet                  []bitmex.Wallet

	Results error
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, l *log.Logger) error {
	p := New(ctx, l)
	defer p.Close()

	// subscribe private
	if p.Auth != nil {
		if err := p.signture(); err != nil {
			return err
		}
	}

	requests, err := p.subscribe(channels, symbols)
	if err != nil {
		// tls: use of closed connection
		p.log.Printf("disconnect %v", err)
		return fmt.Errorf("disconnect %v", err)
	}
	defer p.unsubscribe(requests)

	go p.ping(ctx)

	var eg errgroup.Group
	eg.Go(func() error {
		for {
			p.conn.SetReadDeadline(time.Now().Add(READDEADLINE))
			_, msg, err := p.conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("can't receive error: %v", err)
			}
			// start := time.Now()
			// fmt.Printf("-------------%+v\n", string(msg))

			var r Response

			name, err := jsonparser.GetString(msg, "table")
			if err != nil {
				continue
			}
			r.Action, _ = jsonparser.GetString(msg, "action")
			data, _, _, err := jsonparser.Get(msg, "data")
			if err != nil {
				continue
			}

			switch {
			case strings.HasPrefix(name, "quote"):
				r.Types = Quote
				if err := json.Unmarshal(data, &r.Quote); err != nil {
					continue
				}

			case strings.HasPrefix(name, "orderBookL"):
				r.Types = OrderbookL
				if err := json.Unmarshal(data, &r.OrderbookL); err != nil {
					continue
				}

			case strings.HasPrefix(name, "orderBook"):
				r.Types = Orderbook
				if err := json.Unmarshal(data, &r.Orderbook); err != nil {
					continue
				}

			case strings.HasPrefix(name, "tradeBin"):
				r.Types = TradeBin
				if err := json.Unmarshal(data, &r.TradeBin); err != nil {
					continue
				}

			case strings.HasPrefix(name, "trade"):
				r.Types = Trade
				if err := json.Unmarshal(data, &r.Trade); err != nil {
					continue
				}

			case strings.HasPrefix(name, "announcement"):
				r.Types = Announcement
				if err := json.Unmarshal(data, &r.Announcement); err != nil {
					continue
				}

			case strings.HasPrefix(name, "chat"):
				r.Types = Chat
				if err := json.Unmarshal(data, &r.Chat); err != nil {
					continue
				}

			case strings.HasPrefix(name, "connected"):
				r.Types = Connected
				if err := json.Unmarshal(data, &r.Connected); err != nil {
					continue
				}

			case strings.HasPrefix(name, "funding"):
				r.Types = Funding
				if err := json.Unmarshal(data, &r.Funding); err != nil {
					continue
				}

			case strings.HasPrefix(name, "instrument"):
				r.Types = Instrument
				if err := json.Unmarshal(data, &r.Instrument); err != nil {
					continue
				}

			case strings.HasPrefix(name, "insurance"):
				r.Types = Insurance
				if err := json.Unmarshal(data, &r.Insurance); err != nil {
					continue
				}

			case strings.HasPrefix(name, "settlement"):
				r.Types = Settlement
				if err := json.Unmarshal(data, &r.Settlement); err != nil {
					continue
				}

			case strings.HasPrefix(name, "publicNotifications"):
				r.Types = Notifications
				if err := json.Unmarshal(data, &r.Notifications); err != nil {
					continue
				}

			// Private
			case strings.HasPrefix(name, "execution"):
				r.Types = Execution
				if err := json.Unmarshal(data, &r.Execution); err != nil {
					continue
				}
			case strings.HasPrefix(name, "order"):
				r.Types = Order
				if err := json.Unmarshal(data, &r.Order); err != nil {
					continue
				}
			case strings.HasPrefix(name, "margin"):
				r.Types = Margin
				if err := json.Unmarshal(data, &r.Margin); err != nil {
					continue
				}
			case strings.HasPrefix(name, "position"):
				r.Types = Position
				if err := json.Unmarshal(data, &r.Position); err != nil {
					continue
				}
			case strings.HasPrefix(name, "transact"):
				r.Types = Transact
				if err := json.Unmarshal(data, &r.Transact); err != nil {
					continue
				}
			case strings.HasPrefix(name, "wallet"):
				r.Types = Wallet
				if err := json.Unmarshal(data, &r.Wallet); err != nil {
					continue
				}
			case strings.HasPrefix(name, "privateNotifications"):
				r.Types = NotificationsForPrivate
				if err := json.Unmarshal(data, &r.NotificationsForPrivate); err != nil {
					continue
				}
			case strings.HasPrefix(name, "affiliate"):
				r.Types = Affiliate
				if err := json.Unmarshal(data, &r.Affiliate); err != nil {
					continue
				}

			default:
				r.Types = Undefined
				r.Results = fmt.Errorf("%v", string(msg))
			}

			// switch with ProductCode
			switch {
			case strings.HasSuffix(name, string(bitmex.XBTUSD)):
				r.ProductCode = string(bitmex.XBTUSD)

			case strings.HasSuffix(name, string(bitmex.ETHUSD)):
				r.ProductCode = string(bitmex.ETHUSD)

			case strings.HasSuffix(name, "XRPUSD"):
				r.ProductCode = "XRPUSD"

			// 定義外他通貨
			case strings.Contains(name, ":"):
				r.ProductCode = name[strings.Index(name, ":")+1:]

			default:
				r.ProductCode = "undefined"
			}

			select { // 外部からの停止
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// log.Debugf("recieve to send time: %v\n", time.Now().Sub(start))
			ch <- r
		}
	})

	if err := eg.Wait(); err != nil {
		p.log.Printf("%v", err)

		// 外部からのキャンセル
		if strings.Contains(err.Error(), context.Canceled.Error()) {
			// defer close()/unsubscribe()
			return fmt.Errorf("context stop outside %v", err)
		}
	}

	return nil
}

func (p *Client) subscribe(channels, symbols []string) (requests []Request, err error) {
	var args []interface{}
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				args = append(args, fmt.Sprintf("%s:%s", channels[i], symbols[j]))
			}
		}

		requests = append(requests, Request{
			Op:   "subscribe",
			Args: args,
		})
	} else {
		for i := range channels {
			args = append(args, channels[i])
		}

		requests = append(requests, Request{
			Op:   "subscribe",
			Args: args,
		})
	}

	for i := range requests {
		if err := p.conn.WriteJSON(requests[i]); err != nil {
			return nil, err
		}
		p.log.Printf("subscribed: %v", requests[i])
	}

	return requests, nil
}

func (p *Client) unsubscribe(requests []Request) {
	for i := range requests {
		requests[i].Op = "unsubscribe"
		if err := p.conn.WriteJSON(requests[i]); err != nil {
			_, file, line, _ := runtime.Caller(0)
			p.log.Printf("file:%s, line:%d, error:%+v", file, line+1, err)
		}
	}

	p.log.Println("killed subscribe")
}

func (p *Client) ping(ctx context.Context) {
	var ticker = time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := p.conn.WriteMessage(websocket.TextMessage, []byte(`ping`)); err != nil {
				continue
			}

		case <-ctx.Done():
			return
		}
	}
}

func (p *Client) signture() error {
	// 24時間
	expire := time.Now().UTC().Add(24 * 60 * time.Minute).Unix()
	h := hmac.New(sha256.New, []byte(p.Auth.Secret))
	payload := fmt.Sprintf("GET/realtime%d", expire)
	h.Write([]byte(payload))
	sign := hex.EncodeToString(h.Sum(nil))

	var req = &Request{
		Op: "authKeyExpires",
		ID: 1,
	}
	req.Args = append(req.Args, p.Auth.Key)
	req.Args = append(req.Args, expire)
	req.Args = append(req.Args, sign)

	if err := p.conn.WriteJSON(req); err != nil {
		return err
	}

	return nil
}
