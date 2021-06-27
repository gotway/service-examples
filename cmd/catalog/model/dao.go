package model

import (
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/gotway/service-examples/pkg/catalog"
)

// ProductDAO access to product data
type ProductDAO struct {
	mux      sync.Mutex
	products []catalog.Product
}

// GetProducts obtains products in batches
func (dao *ProductDAO) GetProducts(
	offset int,
	limit int,
) (*catalog.ProductPage, *catalog.ProductError) {
	dao.mux.Lock()
	defer dao.mux.Unlock()
	if dao.products == nil || len(dao.products) == 0 || offset > len(dao.products) {
		return nil, catalog.NotFoundError
	}

	products := make([]catalog.Product, len(dao.products))
	copy(products, dao.products)
	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})

	lowerIndex := offset
	upperIndex := min(offset+limit, len(products))
	slicedProducts := products[lowerIndex:upperIndex]
	if len(slicedProducts) == 0 {
		return nil, catalog.NotFoundError
	}

	productPage := &catalog.ProductPage{slicedProducts, len(products)}
	return productPage, nil
}

// FindProduct finds a product by id
func (dao *ProductDAO) FindProduct(id int) (*catalog.Product, *catalog.ProductError) {
	dao.mux.Lock()
	defer dao.mux.Unlock()
	index := dao.findProductIndex(id)
	if index == -1 {
		return nil, catalog.NotFoundError
	}
	product := dao.products[index]
	return &product, nil
}

// AddProduct adds a product
func (dao *ProductDAO) AddProduct(p *catalog.Product) {
	dao.mux.Lock()
	defer dao.mux.Unlock()
	p.ID = rand.Intn(math.MaxInt32)
	dao.products = append(dao.products, *p)
}

// DeleteProduct deletes a product
func (dao *ProductDAO) DeleteProduct(id int) (bool, *catalog.ProductError) {
	dao.mux.Lock()
	defer dao.mux.Unlock()
	index := dao.findProductIndex(id)
	if index == -1 {
		return false, catalog.NotFoundError
	}
	dao.products = append(dao.products[:index], dao.products[index+1:]...)
	return true, nil
}

// UpdateProduct updates a product
func (dao *ProductDAO) UpdateProduct(id int, p *catalog.Product) (bool, *catalog.ProductError) {
	dao.mux.Lock()
	defer dao.mux.Unlock()
	index := dao.findProductIndex(id)
	if index == -1 {
		return false, catalog.NotFoundError
	}
	p.ID = id
	dao.products[index] = *p
	return true, nil
}

func (dao *ProductDAO) findProductIndex(id int) int {
	for index, p := range dao.products {
		if p.ID == id {
			return index
		}
	}
	return -1
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
