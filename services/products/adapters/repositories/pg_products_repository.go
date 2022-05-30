package repositories

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/lucasmls/ecommerce/services/products/adapters/repositories/models"
	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PgProductsRepository struct {
	db *sql.DB
}

func NewPgProductsRepository(connectionString string) (*PgProductsRepository, error) {
	db, err := sql.Open(
		"postgres",
		connectionString,
	)
	if err != nil {
		return nil, err
	}

	repository := &PgProductsRepository{
		db: db,
	}

	return repository, nil
}

func MustNewPgProductsRepository(connectionString string) *PgProductsRepository {
	repo, err := NewPgProductsRepository(connectionString)
	if err != nil {
		panic(err)
	}

	return repo
}

func (r *PgProductsRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	p := models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	err := p.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return domain.Product{}, err
	}

	response := domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	return response, nil
}

func (r *PgProductsRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	p, err := models.FindProduct(ctx, r.db, product.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, domain.ErrProductNotFound
		}

		return domain.Product{}, err
	}

	p.Name = product.Name
	p.Description = product.Description
	p.Price = product.Price

	_, err = p.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return domain.Product{}, err
	}

	response := domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	return response, nil
}

func (r *PgProductsRepository) Delete(ctx context.Context, id int) error {
	product, err := models.FindProduct(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrProductNotFound
		}

		return err
	}

	_, err = product.Delete(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgProductsRepository) List(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, error) {
	products, err := models.Products().All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	var response []domain.Product
	for _, product := range products {
		response = append(response, domain.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return response, err
}
