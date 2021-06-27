package service

import (
	"github.com/gotway/service-examples/cmd/catalog/model"
	"github.com/gotway/service-examples/pkg/catalog"
)

// ProductService manages product business logic
type ProductService struct {
	dao model.ProductDAO
}

// GetProducts obtains products in batches
func (s *ProductService) GetProducts(
	offset int,
	limit int,
) (*catalog.ProductPage, *catalog.ProductError) {
	return s.dao.GetProducts(offset, limit)
}

// FindProduct finds a product by id
func (s *ProductService) FindProduct(id int) (*catalog.Product, *catalog.ProductError) {
	return s.dao.FindProduct(id)
}

// AddProduct adds a product
func (s *ProductService) AddProduct(p *catalog.Product) {
	s.dao.AddProduct(p)
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id int) (bool, *catalog.ProductError) {
	return s.dao.DeleteProduct(id)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(id int, p *catalog.Product) (bool, *catalog.ProductError) {
	return s.dao.UpdateProduct(id, p)
}
