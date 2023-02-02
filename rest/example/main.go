package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/client"
	"github.com/polygon-io/client-go/rest/iter"
	"github.com/polygon-io/client-go/rest/models"
)

func main() {
	// call API for aggregates
	iter := getAggregates("X:BTCUSD")

	// loop through and print aggregates
	for iter.Next() {

		if iter.Err() != nil {
			log.Fatal(iter.Err())
		}

		aggMap := map[string]any{
			// "T":   iter.Item().Ticker, not getting value back
			"c":   iter.Item().Close,
			"h":   iter.Item().High,
			"l":   iter.Item().Low,
			"n":   iter.Item().Transactions,
			"o":   iter.Item().Open,
			"v":   iter.Item().Volume,
			"vw":  iter.Item().VWAP,
			"otc": iter.Item().OTC,
		}
		log.Print(aggMap)
	}
}

// returns iterator to aggregates
func getAggregates(ticker string) *iter.Iter[models.Agg] {

	// create a client
	c := polygon.AggsClient{
		Client: client.New(getEnvVar("POLYGON_API_KEY")),
	}

	// call api with parameters
	iter := c.ListAggs(context.Background(), models.ListAggsParams{
		Ticker:     ticker,
		Multiplier: 1,
		Timespan:   "day",
		From:       models.Millis(time.Date(2021, 7, 22, 0, 0, 0, 0, time.UTC)),
		To:         models.Millis(time.Date(2021, 8, 22, 0, 0, 0, 0, time.UTC)),
	}.WithOrder(models.Desc).WithLimit(2).WithAdjusted(true))
	
	return iter
}

// returns the value of the key from env file
func getEnvVar(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
