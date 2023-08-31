package opensea

import (
	"github.com/foundVanting/opensea-stream-go/opensea"
	"math/big"
)

const (
	moduleName = "OpenSea"
	host       = ""
)

type Client struct {
	ApiKey       string
	StreamClient *opensea.StreamClient
}

type Sale struct {
	Collection     string
	CollectionLink string
	Timestamp      string
	Item           string
	Seller         string
	Buyer          string
	Image          string
	Username       string
	PriceInfo      PriceConversion
}

type Listing struct {
	Collection     string
	CollectionLink string
	Timestamp      int64
	Item           string
	ItemURL        string
	Image          string
	WeiPrice       *big.Float
	EthereumPrice  float64
	Price          int64
	Balance        string
	Seller         string
	SellerLink     string
	Dex            string
	Symbol         string
	PriceInfo      PriceConversion
}

type PriceConversion struct {
	PriceBefore  float64
	Price        *big.Float
	Floor        float64
	PercentMinus float64
	Difference   string
}

type CollectionData struct {
	Stats struct {
		OneHourVolume         int     `json:"one_hour_volume"`
		OneHourChange         int     `json:"one_hour_change"`
		OneHourSales          int     `json:"one_hour_sales"`
		OneHourSalesChange    int     `json:"one_hour_sales_change"`
		OneHourAveragePrice   int     `json:"one_hour_average_price"`
		OneHourDifference     int     `json:"one_hour_difference"`
		SixHourVolume         float64 `json:"six_hour_volume"`
		SixHourChange         int     `json:"six_hour_change"`
		SixHourSales          int     `json:"six_hour_sales"`
		SixHourSalesChange    int     `json:"six_hour_sales_change"`
		SixHourAveragePrice   float64 `json:"six_hour_average_price"`
		SixHourDifference     float64 `json:"six_hour_difference"`
		OneDayVolume          float64 `json:"one_day_volume"`
		OneDayChange          float64 `json:"one_day_change"`
		OneDaySales           int     `json:"one_day_sales"`
		OneDaySalesChange     float64 `json:"one_day_sales_change"`
		OneDayAveragePrice    float64 `json:"one_day_average_price"`
		OneDayDifference      float64 `json:"one_day_difference"`
		SevenDayVolume        float64 `json:"seven_day_volume"`
		SevenDayChange        float64 `json:"seven_day_change"`
		SevenDaySales         int     `json:"seven_day_sales"`
		SevenDayAveragePrice  float64 `json:"seven_day_average_price"`
		SevenDayDifference    float64 `json:"seven_day_difference"`
		ThirtyDayVolume       float64 `json:"thirty_day_volume"`
		ThirtyDayChange       float64 `json:"thirty_day_change"`
		ThirtyDaySales        int     `json:"thirty_day_sales"`
		ThirtyDayAveragePrice float64 `json:"thirty_day_average_price"`
		ThirtyDayDifference   float64 `json:"thirty_day_difference"`
		TotalVolume           float64 `json:"total_volume"`
		TotalSales            int     `json:"total_sales"`
		TotalSupply           int     `json:"total_supply"`
		Count                 int     `json:"count"`
		NumOwners             int     `json:"num_owners"`
		AveragePrice          float64 `json:"average_price"`
		NumReports            int     `json:"num_reports"`
		MarketCap             float64 `json:"market_cap"`
		FloorPrice            float64 `json:"floor_price"`
	} `json:"stats"`
}
