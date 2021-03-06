package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gotway/service-examples/cmd/stock/service"
	"github.com/gotway/service-examples/pkg/stock"

	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func upsertStock(w http.ResponseWriter, r *http.Request) {
	productID, err := getProductID(r)
	if err != nil {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var stockData stock.StockData
	_ = json.NewDecoder(r.Body).Decode(&stockData)
	if !stockData.IsValid() {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stock := stockData.ToStock(productID)
	resultStock, stockError := service.UpsertStock(productID, &stock)
	if stockError != nil {
		handleError(w, stockError)
		return
	}
	setHeaders(w)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resultStock)
}

func upsertStockList(w http.ResponseWriter, r *http.Request) {
	var stockList stock.StockList
	_ = json.NewDecoder(r.Body).Decode(&stockList)
	if !stockList.IsValid() {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resultStockList := service.UpsertStockList(stockList.Stock)
	if !resultStockList.HasStock() {
		setHeaders(w)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	setHeaders(w)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&resultStockList)
}

func getStock(w http.ResponseWriter, r *http.Request) {
	productID, err := getProductID(r)
	if err != nil {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stock, stockErr := service.GetStock(productID)
	if stockErr != nil {
		handleError(w, stockErr)
		return
	}
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stock)
}

func getStockList(w http.ResponseWriter, r *http.Request) {
	productIDs := getProductIDs(r)
	if productIDs == nil {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stockList := service.GetStockList(productIDs)
	if !stockList.HasStock() {
		setHeaders(w)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&stockList)
}

func getProductID(r *http.Request) (int, error) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func getProductIDs(r *http.Request) []int {
	q := r.URL.Query()
	requestProductIDs := q["productId"]
	if requestProductIDs == nil {
		return nil
	}
	productIDs := []int{}
	for _, val := range requestProductIDs {
		if val == "" {
			continue
		}
		productID, err := strconv.Atoi(val)
		if err == nil {
			productIDs = append(productIDs, productID)
		}
	}
	if len(productIDs) == 0 {
		return nil
	}
	return productIDs
}

func handleError(w http.ResponseWriter, err *stock.StockError) {
	log.Print(err)
	setHeaders(w)
	w.WriteHeader(err.Code)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "s-maxage=0")
}
