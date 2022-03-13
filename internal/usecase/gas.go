package usecase

import (
	"encoding/json"
	"github.com/aresprotocols/trojan-box/internal/app"
	gocache "github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type GasUseCase struct {
	goCache *gocache.Cache
}

const GasPriceKey = "GAS_PRICE_KEY"
const GatePricePrefixKey = "GATE_PRICE_"

func newGas(svc *useCase) *GasUseCase {
	return &GasUseCase{svc.goCache}
}

func (u *GasUseCase) CalGasFeeToAres(gas int64, config app.Config) decimal.Decimal {
	gasPrice := decimal.NewFromFloat(u.fetchGasPrice(config))
	ethGasFee := decimal.NewFromInt(gas).Mul(gasPrice).Div(decimal.New(1, 9))
	aresGasFee := ethGasFee.Div(u.fetchAresEthRatio())
	return aresGasFee
}

func (u *GasUseCase) fetchGasPrice(config app.Config) float64 {
	gasPrice, found := u.goCache.Get(GasPriceKey)
	if found {
		return gasPrice.(float64)
	}
	resp, err := http.Get(config.Ares.GasStationUrl)
	if err != nil {
		logger.WithError(err).Errorf("get gas price from gasstaion occur err")
		return 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Errorf("get body occur error")
		return 0
	}
	logger.Infof("fetch gas price from gassstaion:%s", string(body))
	result := make(map[string]interface{})

	err = json.Unmarshal(body, &result)

	average := result["average"].(float64)
	average = average / 10
	err = u.goCache.Add(GasPriceKey, average, time.Minute)
	if err != nil {
		logger.WithError(err).Errorf("gocache add cache occur err")
	}
	return average
}
func (u *GasUseCase) fetchAresEthRatio() decimal.Decimal {
	aresPrice := u.getPriceFromGate("ares_usdt")
	ethPrice := u.getPriceFromGate("eth_usdt")
	if aresPrice == 0 || ethPrice == 0 {
		return decimal.Zero
	}
	return decimal.NewFromFloat(aresPrice).Div(decimal.NewFromFloat(ethPrice))
}

func (u *GasUseCase) getPriceFromGate(pair string) float64 {
	cachePrice, found := u.goCache.Get(GatePricePrefixKey + pair)
	if found {
		return cachePrice.(float64)
	}
	log := logger.WithField("pari", pair)
	resp, err := http.Get("https://data.gateapi.io/api2/1/ticker/" + pair)
	if err != nil {
		log.WithError(err).Errorf("get gas from gate occur err")
		return 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("get body occur error")
		return 0
	}
	logger.Infof("fetch %s price from gassstaion:%s", pair, string(body))
	result := make(map[string]interface{})

	err = json.Unmarshal(body, &result)
	priceStr := result["highestBid"].(string)
	price, err := strconv.ParseFloat(priceStr, 10)
	if err != nil {
		logger.WithError(err).Errorf("parse price occur err")
		return 0
	}
	err = u.goCache.Add(GatePricePrefixKey+pair, price, time.Minute)
	if err != nil {
		logger.WithError(err).Errorf("gocache add cache occur err")
	}
	return price
}
