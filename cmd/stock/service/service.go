package service

import (
	"strconv"
	"sync"

	"github.com/gotway/service-examples/cmd/stock/redis"
	"github.com/gotway/service-examples/pkg/stock"
)

// UpsertStock upserts the stock of a product
func UpsertStock(productID int, s *stock.Stock) (*stock.Stock, *stock.StockError) {
	key := getKey(productID)
	ttl, err := redis.Set(key, s.Units, s.TTL)
	if err != nil {
		return nil, stock.InternalError(productID, err)
	}
	resultStock := &stock.Stock{ProductID: productID, Units: s.Units, TTL: ttl}
	return resultStock, nil
}

// UpsertStockList upserts the stock of a list of products
func UpsertStockList(stockList []stock.Stock) stock.StockList {
	var wg sync.WaitGroup
	var sl stock.StockList

	wg.Add(len(stockList))
	for _, item := range stockList {
		go func(s stock.Stock) {
			defer wg.Done()
			resultStock, err := UpsertStock(s.ProductID, &s)
			if err == nil {
				sl.AddStock(resultStock)
			} else {
				sl.HandleError(err)
			}
		}(item)
	}
	wg.Wait()

	return sl.Get()
}

// GetStock gets the stock of a product
func GetStock(productID int) (*stock.Stock, *stock.StockError) {
	unitsChan := make(chan stockResult)
	ttlChan := make(chan stockResult)

	go getStockUnits(productID, unitsChan)
	go getTTL(productID, ttlChan)

	units, ttl := <-unitsChan, <-ttlChan
	if units.err != nil {
		return nil, units.err
	}
	if ttl.err != nil {
		return nil, units.err
	}

	stock := &stock.Stock{ProductID: productID, Units: units.val, TTL: ttl.val}
	return stock, nil
}

// GetStockList gets the stock of a list of products
func GetStockList(productIDs []int) stock.StockList {
	var wg sync.WaitGroup
	var sl stock.StockList

	wg.Add(len(productIDs))
	for _, productID := range productIDs {
		go func(id int) {
			defer wg.Done()
			stock, err := GetStock(id)
			if err == nil {
				sl.AddStock(stock)
			} else {
				sl.HandleError(err)
			}
		}(productID)
	}
	wg.Wait()

	return sl.Get()
}

func getKey(productID int) string {
	return strconv.Itoa(productID)
}

type stockResult struct {
	val int
	err *stock.StockError
}

func getStockUnits(productID int, c chan stockResult) {
	key := getKey(productID)
	defer close(c)
	val, err := redis.Get(key)
	if err != nil {
		c <- stockResult{0, stock.OutOfStockError(productID)}
		return
	}
	units, parseErr := strconv.Atoi(val)
	if parseErr != nil {
		c <- stockResult{0, stock.InternalError(productID, parseErr)}
		return
	}
	if units <= 0 {
		c <- stockResult{0, stock.OutOfStockError(productID)}
		return
	}
	c <- stockResult{units, nil}
}

func getTTL(productID int, c chan stockResult) {
	key := getKey(productID)
	defer close(c)
	ttl, err := redis.TTL(key)
	if err != nil {
		c <- stockResult{0, stock.OutOfStockError(productID)}
		return
	}
	c <- stockResult{ttl, nil}
}
