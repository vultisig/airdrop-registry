package volume

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTCVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"actions": []map[string]any{
				{
					"date":   "1739518746459473723",
					"height": "19862079",
					"in": []map[string]any{
						{
							"address": "0x060c27cd6719477f233e403d74da9513886f0a1a",
							"coins": []map[string]any{
								{
									"amount": "315789473676",
									"asset":  "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
								},
							},
							"txID": "6A31B1ADC6211047175A1A465985B7B2E6C28945C587F91305B73946491226CF",
						},
					},
					"metadata": map[string]any{
						"swap": map[string]any{
							"affiliateAddress": "t",
							"affiliateFee":     "0",
							"inPriceUSD":       "1.0005046063264105",
							"isStreamingSwap":  true,
							"liquidityFee":     "349228438",
							"memo":             "=:ETH.THOR-0XA5F2211B9B8170F694421F2046281775E8468044:0x060c27cd6719477F233E403d74da9513886F0a1A:7508256142234/1/0:t:0",
							"networkFees": []map[string]any{
								{
									"amount": "10387967866",
									"asset":  "ETH.THOR-0XA5F2211B9B8170F694421F2046281775E8468044",
								},
								{
									"amount": "272164500",
									"asset":  "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
								},
							},
							"outPriceUSD": "0.06582207333142955",
							"streamingSwapMeta": map[string]any{
								"count": "19",
								"depositedCoin": map[string]any{
									"amount": "500000000000",
									"asset":  "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
								},
								"failedSwapReasons": []string{
									"emit asset 392492334429 less than price limit 392951609490",
									"emit asset 392492334428 less than price limit 392951609490",
									"emit asset 392492334428 less than price limit 392951609490",
									"emit asset 392492334428 less than price limit 392951609490",
									"emit asset 392492334427 less than price limit 392951609490",
									"emit asset 392492334427 less than price limit 392951609490",
									"emit asset 2720179808079 less than price limit 2750661266623",
								},
								"failedSwaps": []string{
									"13", "14", "15", "16", "17", "18", "19",
								},
								"inCoin": map[string]any{
									"amount": "315789473676",
									"asset":  "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
								},
								"interval":   "1",
								"lastHeight": "19862097",
								"outCoin": map[string]any{
									"amount": "4757594875611",
									"asset":  "ETH.THOR-0XA5F2211B9B8170F694421F2046281775E8468044",
								},
								"quantity": "19",
							},
							"swapSlip":   "16",
							"swapTarget": "395171375897",
							"txType":     "swap",
						},
					},
					"out": []map[string]any{
						{
							"address":   "thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd",
							"affiliate": true,
							"coins": []map[string]any{
								{
									"amount": "996806325",
									"asset":  "THOR.RUNE",
								},
							},
							"height": "19738122",
							"txID":   "",
						},
						{
							"address": "0x060c27cd6719477f233e403d74da9513886f0a1a",
							"coins": []map[string]any{
								{
									"amount": "4747206907745",
									"asset":  "ETH.THOR-0XA5F2211B9B8170F694421F2046281775E8468044",
								},
							},
							"height": "19862103",
							"txID":   "2FFFCAF7A20E138A2558B2F21D29A227D7EA4D0B806C588CA81AA07CEDF3FE9B",
						},
						{
							"address": "0x060c27cd6719477f233e403d74da9513886f0a1a",
							"coins": []map[string]any{
								{
									"amount": "183938361800",
									"asset":  "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
								},
							},
							"height": "19862103",
							"txID":   "A28AF45D9116FAED319DC78AEBEDA731533EC2D9AB00602E7F240A015BEA1761",
						},
					},
					"pools": []string{
						"ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
						"ETH.THOR-0XA5F2211B9B8170F694421F2046281775E8468044",
					},
					"status": "success",
					"type":   "swap",
				},
			},
			"meta": map[string]any{
				"nextPageToken": "198620799000000001",
				"prevPageToken": "198620799000000001",
			},
		}

		response2 := map[string]any{
			"actions": []any{},
			"meta": map[string]any{
				"nextPageToken": "",
				"prevPageToken": "",
			},
		}

		if r.URL.Query().Get("nextPageToken") == "198620799000000001" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response2)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer mockServer.Close()
	vr := tcVolumeTrack{
		logger:  logrus.WithField("module", "tc_vol_service").Logger,
		baseUrl: mockServer.URL,
	}
	expect := map[string]float64{
		"0x060c27cd6719477f233e403d74da9513886f0a1a": 324578205539.9229,
	}
	res, err := vr.processVolume(1739510000, 1739519656, "t")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, expect, res)
}
