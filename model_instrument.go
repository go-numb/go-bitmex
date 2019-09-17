package bitmex

import (
	"time"
)

// Tradeable Contracts, Indices, and History
type Instrument struct {
	Symbol                         string    `json:"symbol"`
	RootSymbol                     string    `json:"rootSymbol,omitempty"`
	State                          string    `json:"state,omitempty"`
	Typ                            string    `json:"typ,omitempty"`
	Listing                        time.Time `json:"listing,omitempty"`
	Front                          time.Time `json:"front,omitempty"`
	Expiry                         time.Time `json:"expiry,omitempty"`
	Settle                         time.Time `json:"settle,omitempty"`
	RelistInterval                 time.Time `json:"relistInterval,omitempty"`
	InverseLeg                     string    `json:"inverseLeg,omitempty"`
	SellLeg                        string    `json:"sellLeg,omitempty"`
	BuyLeg                         string    `json:"buyLeg,omitempty"`
	OptionStrikePcnt               float64   `json:"optionStrikePcnt,omitempty"`
	OptionStrikeRound              float64   `json:"optionStrikeRound,omitempty"`
	OptionStrikePrice              float64   `json:"optionStrikePrice,omitempty"`
	OptionMultiplier               float64   `json:"optionMultiplier,omitempty"`
	PositionCurrency               string    `json:"positionCurrency,omitempty"`
	Underlying                     string    `json:"underlying,omitempty"`
	QuoteCurrency                  string    `json:"quoteCurrency,omitempty"`
	UnderlyingSymbol               string    `json:"underlyingSymbol,omitempty"`
	Reference                      string    `json:"reference,omitempty"`
	ReferenceSymbol                string    `json:"referenceSymbol,omitempty"`
	CalcInterval                   time.Time `json:"calcInterval,omitempty"`
	PublishInterval                time.Time `json:"publishInterval,omitempty"`
	PublishTime                    time.Time `json:"publishTime,omitempty"`
	MaxOrderQty                    int       `json:"maxOrderQty,omitempty"`
	MaxPrice                       float64   `json:"maxPrice,omitempty"`
	LotSize                        int       `json:"lotSize,omitempty"`
	TickSize                       float64   `json:"tickSize,omitempty"`
	Multiplier                     int       `json:"multiplier,omitempty"`
	SettlCurrency                  string    `json:"settlCurrency,omitempty"`
	UnderlyingToPositionMultiplier int       `json:"underlyingToPositionMultiplier,omitempty"`
	UnderlyingToSettleMultiplier   int       `json:"underlyingToSettleMultiplier,omitempty"`
	QuoteToSettleMultiplier        int       `json:"quoteToSettleMultiplier,omitempty"`
	IsQuanto                       bool      `json:"isQuanto,omitempty"`
	IsInverse                      bool      `json:"isInverse,omitempty"`
	InitMargin                     float64   `json:"initMargin,omitempty"`
	MaintMargin                    float64   `json:"maintMargin,omitempty"`
	RiskLimit                      int       `json:"riskLimit,omitempty"`
	RiskStep                       int       `json:"riskStep,omitempty"`
	Limit                          float64   `json:"limit,omitempty"`
	Capped                         bool      `json:"capped,omitempty"`
	Taxed                          bool      `json:"taxed,omitempty"`
	Deleverage                     bool      `json:"deleverage,omitempty"`
	MakerFee                       float64   `json:"makerFee,omitempty"`
	TakerFee                       float64   `json:"takerFee,omitempty"`
	SettlementFee                  float64   `json:"settlementFee,omitempty"`
	InsuranceFee                   float64   `json:"insuranceFee,omitempty"`
	FundingBaseSymbol              string    `json:"fundingBaseSymbol,omitempty"`
	FundingQuoteSymbol             string    `json:"fundingQuoteSymbol,omitempty"`
	FundingPremiumSymbol           string    `json:"fundingPremiumSymbol,omitempty"`
	FundingTimestamp               time.Time `json:"fundingTimestamp,omitempty"`
	FundingInterval                time.Time `json:"fundingInterval,omitempty"`
	FundingRate                    float64   `json:"fundingRate,omitempty"`
	IndicativeFundingRate          float64   `json:"indicativeFundingRate,omitempty"`
	RebalanceTimestamp             time.Time `json:"rebalanceTimestamp,omitempty"`
	RebalanceInterval              time.Time `json:"rebalanceInterval,omitempty"`
	OpeningTimestamp               time.Time `json:"openingTimestamp,omitempty"`
	ClosingTimestamp               time.Time `json:"closingTimestamp,omitempty"`
	SessionInterval                time.Time `json:"sessionInterval,omitempty"`
	PrevClosePrice                 float64   `json:"prevClosePrice,omitempty"`
	LimitDownPrice                 float64   `json:"limitDownPrice,omitempty"`
	LimitUpPrice                   float64   `json:"limitUpPrice,omitempty"`
	BankruptLimitDownPrice         float64   `json:"bankruptLimitDownPrice,omitempty"`
	BankruptLimitUpPrice           float64   `json:"bankruptLimitUpPrice,omitempty"`
	PrevTotalVolume                int       `json:"prevTotalVolume,omitempty"`
	TotalVolume                    int       `json:"totalVolume,omitempty"`
	Volume                         int       `json:"volume,omitempty"`
	Volume24h                      int       `json:"volume24h,omitempty"`
	PrevTotalTurnover              int       `json:"prevTotalTurnover,omitempty"`
	TotalTurnover                  int       `json:"totalTurnover,omitempty"`
	Turnover                       int       `json:"turnover,omitempty"`
	Turnover24h                    int       `json:"turnover24h,omitempty"`
	PrevPrice24h                   float64   `json:"prevPrice24h,omitempty"`
	Vwap                           float64   `json:"vwap,omitempty"`
	HighPrice                      float64   `json:"highPrice,omitempty"`
	LowPrice                       float64   `json:"lowPrice,omitempty"`
	LastPrice                      float64   `json:"lastPrice,omitempty"`
	LastPriceProtected             float64   `json:"lastPriceProtected,omitempty"`
	LastTickDirection              string    `json:"lastTickDirection,omitempty"`
	LastChangePcnt                 float64   `json:"lastChangePcnt,omitempty"`
	BidPrice                       float64   `json:"bidPrice,omitempty"`
	MidPrice                       float64   `json:"midPrice,omitempty"`
	AskPrice                       float64   `json:"askPrice,omitempty"`
	ImpactBidPrice                 float64   `json:"impactBidPrice,omitempty"`
	ImpactMidPrice                 float64   `json:"impactMidPrice,omitempty"`
	ImpactAskPrice                 float64   `json:"impactAskPrice,omitempty"`
	HasLiquidity                   bool      `json:"hasLiquidity,omitempty"`
	OpenInterest                   int       `json:"openInterest,omitempty"`
	OpenValue                      int       `json:"openValue,omitempty"`
	FairMethod                     string    `json:"fairMethod,omitempty"`
	FairBasisRate                  float64   `json:"fairBasisRate,omitempty"`
	FairBasis                      float64   `json:"fairBasis,omitempty"`
	FairPrice                      float64   `json:"fairPrice,omitempty"`
	MarkMethod                     string    `json:"markMethod,omitempty"`
	MarkPrice                      float64   `json:"markPrice,omitempty"`
	IndicativeTaxRate              float64   `json:"indicativeTaxRate,omitempty"`
	IndicativeSettlePrice          float64   `json:"indicativeSettlePrice,omitempty"`
	OptionUnderlyingPrice          float64   `json:"optionUnderlyingPrice,omitempty"`
	SettledPrice                   float64   `json:"settledPrice,omitempty"`
	Timestamp                      time.Time `json:"timestamp,omitempty"`
}