package restful

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/maurodelazeri/bitmex-go/swagger"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type OrderApi struct {
	swaggerOrderApi *swagger.OrderApiService
	ctx             context.Context
}

func NewOrderApi(swaggerOrderApi *swagger.OrderApiService, ctx context.Context) *OrderApi {
	return &OrderApi{swaggerOrderApi: swaggerOrderApi, ctx: ctx}
}

// BUY
func (o *OrderApi) LimitBuy(symbol string, orderQty int, price float64, clientOrderIDPrefix string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	if price <= 0 {
		return nil, swagger.Order{}, errors.New("price must be positive")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"ordType":  "Limit",
		"orderQty": float32(orderQty),
		"price":    price,
		"clOrdID":  clOrdID,
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}
	return response, order, nil
}

func (o *OrderApi) MarketBuy(symbol string, orderQty int, clientOrderIDPrefix string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"ordType":  "Market",
		"orderQty": float32(orderQty),
		"clOrdID":  clOrdID,
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}
	return response, order, nil
}

// SELL
func (o *OrderApi) LimitSell(symbol string, orderQty int, price float64, clientOrderIDPrefix string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	if price <= 0 {
		return nil, swagger.Order{}, errors.New("price must be positive")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"orderQty": float32(-orderQty),
		"price":    price,
		"clOrdID":  clOrdID,
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}
	return response, order, nil
}

func (o *OrderApi) MarketSell(symbol string, orderQty int, clientOrderIDPrefix string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"orderQty": float32(-orderQty),
		"clOrdID":  clOrdID,
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}
	return response, order, nil
}

// StopOrder: Like a Stop Market, but enters a Limit order instead of a Market order. Specify an orderQty, stopPx, and price.
func (o *OrderApi) StopOrder(symbol string, orderQty int, price float64, stopPx float64, clientOrderIDPrefix string, positionSide string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	if price <= 0 {
		return nil, swagger.Order{}, errors.New("price must be positive")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	var amount float32

	if positionSide == "Sell" {
		amount = float32(-orderQty)
	} else {
		amount = float32(orderQty)
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"orderQty": amount,
		//"price":    price,
		"stopPx":  stopPx,
		"clOrdID": clOrdID,
		"ordType": "Stop",
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}
	return response, order, nil
}

// StopLimitAmend change a current order in place
func (o *OrderApi) SetAmendOrder(symbol string, orderQty int, price float64, stopPx float64, OrderID string) (resp *http.Response, order swagger.Order, err error) {
	if symbol == "" {
		return nil, swagger.Order{}, errors.New("symbol can NOT be empty")
	}
	if OrderID == "" {
		return nil, swagger.Order{}, errors.New("oderID can NOT be empty")
	}
	if price <= 0 {
		return nil, swagger.Order{}, errors.New("price must be positive")
	}

	params := map[string]interface{}{
		"price":   price,
		"stopPx":  stopPx,
		"orderID": OrderID,
	}

	order, response, err := o.swaggerOrderApi.OrderAmend(o.ctx, params)

	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}

	return response, order, nil
}

// TakeProfit
func (o *OrderApi) TakeProfit(symbol string, orderQty float64, price float64, stopPx float64, clientOrderIDPrefix string, positionSide string) (resp *http.Response, orderId string, err error) {
	if symbol == "" {
		return nil, "", errors.New("symbol can NOT be empty")
	}
	if price <= 0 {
		return nil, "", errors.New("price must be positive")
	}
	clOrdID := ""
	if clientOrderIDPrefix != "" {
		s := strings.Replace(base64.StdEncoding.EncodeToString(uuid.NewV4().Bytes()), "=", "", -1)
		clOrdID = clientOrderIDPrefix + s
	}

	var amount float32

	if positionSide == "Buy" {
		amount = float32(-orderQty)
	} else {
		amount = float32(orderQty)
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"orderQty": amount,
		"price":    price,
		"stopPx":   stopPx,
		"clOrdID":  clOrdID,
		"ordType":  "StopLimit",
	}
	order, response, err := o.swaggerOrderApi.OrderNew(o.ctx, symbol, params)
	if err != nil || response.StatusCode != 200 {
		return response, order.OrderID, err
	}
	return response, order.OrderID, nil
}

// GetAllordersHistory returns the orders history
func (o *OrderApi) GetAllordersHistory(symbol string, filter string, count float32, reverse bool, startTime time.Time, endTime time.Time) (resp *http.Response, order []swagger.Order, err error) {

	var orderHistory []swagger.Order

	if symbol == "" {
		return nil, orderHistory, errors.New("symbol can NOT be empty")
	}

	params := map[string]interface{}{
		"symbol":    symbol,
		"count":     count,
		"reverse":   reverse,
		"startTime": startTime,
		"endTime":   endTime,
	}

	orders, response, err := o.swaggerOrderApi.OrderGetOrders(o.ctx, params)
	if err != nil || response.StatusCode != 200 {
		return response, order, err
	}

	return response, orders, nil
}
