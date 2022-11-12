package service

import (
	"fmt"

	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/features/currencies/domain/values"

	// TODO: replace this with a proper facade
	SDK "github.com/CoinAPI/coinapi-sdk/data-api/go-rest/v1"
) 

type ExchangeRatesGetter = func(currencies values.Currencies) (values.ExchangeRates, error)

const baseAsset = "USD"

func contains(s []string, e string) bool {
  for _, sliceElem := range s {
    if sliceElem == e {
      return true
    }
  }
  return false
}

func NewExchangeRatesGetter(sdk *SDK.SDK) ExchangeRatesGetter {
  return func(currencies values.Currencies) (values.ExchangeRates, error) {
    allRates, err := sdk.Exchange_rates_get_all_current_rates(baseAsset)
    if err != nil {
      return values.ExchangeRates{}, core_err.Rethrow("getting all exchange rates from an external api", err)
    }
    rates := values.ExchangeRates{}
    for _, rate := range allRates {
      if contains(currencies, rate.Asset_id_quote) {
        rates[rate.Asset_id_quote] = rate.Rate.InexactFloat64()
      }
    }
    if len(rates) != len(currencies) {
      return rates, fmt.Errorf("the external api provided only %+v rates, but %+v currencies were requested", rates, currencies)
    }
    return rates, nil
  }
}


