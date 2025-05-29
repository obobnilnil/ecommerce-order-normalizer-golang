package repository

import (
	"database/sql"
	"e-commerce/internal/e-commerce/model"
)

type RepositoryPort interface {
	NormalizeOrderRepository(orders []model.CleanedOrder) error
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}
func (r *repositoryAdapter) NormalizeOrderRepository(orders []model.CleanedOrder) error {
	return nil
}
