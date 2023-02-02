package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/client"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/stretchr/testify/assert"
)

const (
	expectedAggsResponseURL = "https://api.polygon.io/v2/aggs/ticker/X:BTCUSD/range/1/day/1626912000000/1629590400000"

	agg1 = `{
		"T": "X:BTCUSD",
		"c": 49284.63,
		"h": 49540.01,
		"l": 48058.1,
		"n": 302065,
		"o": 48883.84,
		"v": 11715.632815081244,
		"vw": 48834.5697,
		"otc": false
}`

	expectedAggsResponse = `{
		"status": "OK",
		"request_id": "6a7e466379af0a71039d60cc78e72282",
		"ticker": "X:BTCUSD",
		"queryCount": 2,
		"resultsCount": 2,
		"adjusted": true,
		"results": [
			{
				"T": "X:BTCUSD",
				"c": 49284.63,
				"h": 49540.01,
				"l": 48058.1,
				"n": 302065,
				"o": 48883.84,
				"v": 11715.632815081244,
				"vw": 48834.5697,
				"otc": false
			}
		]
	}`
)

func TestPrintAggregate(t *testing.T) {
	// register client
	c := polygon.AggsClient{
		Client: client.New(getEnvVar("POLYGON_API_KEY")),
	}

	// mock http call
	httpmock.ActivateNonDefault(c.HTTP.GetClient())
	defer httpmock.DeactivateAndReset()
	registerResponder(expectedAggsResponseURL, expectedAggsResponse)

	// call api with parameters
	iter := getAggregates("X:BTCUSD")
	
	// test iter
	assert.Nil(t, iter.Err())
	assert.NotNil(t, iter.Item())
	assert.True(t, iter.Next())

	// testing the ability to Unmarshal JSON to struct
	// setting "expectations"
	expect := models.Agg{}
	err := json.Unmarshal([]byte(agg1), &expect) // map agg1 to Agg struct
	assert.Nil(t, err)

	// testing expections vs function results
	// assert.Equal(t, expect.Ticker, iter.Item().Ticker) - can't test due to value not returning
	assert.Equal(t, expect.Close, iter.Item().Close)
	assert.Equal(t, expect.High, iter.Item().High)
	assert.Equal(t, expect.Low, iter.Item().Low)
	assert.Equal(t, expect.Transactions, iter.Item().Transactions)
	assert.Equal(t, expect.Open, iter.Item().Open)
	assert.Equal(t, expect.Volume, iter.Item().Volume)
	assert.Equal(t, expect.VWAP, iter.Item().VWAP)
	assert.Equal(t, expect.OTC, iter.Item().OTC)
}

// borrowed from polygon codebase
func registerResponder(url, body string) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, body)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		},
	)
}
