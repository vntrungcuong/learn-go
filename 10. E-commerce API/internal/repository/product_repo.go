package repository

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/models"
)

func GetProducts(limit int) ([]models.Product, error) {
	// Query trực tiếp hiệu năng cao
	rows, err := config.DB.Query(context.Background(),
		"SELECT id, category_id, name, price, stock, description FROM products LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Stock, &p.Description); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
