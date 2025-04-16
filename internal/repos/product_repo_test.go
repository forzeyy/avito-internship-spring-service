package repos_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewProductRepo(mock)

	product := &models.Product{
		ID:          uuid.New(),
		Type:        "box",
		ReceptionID: uuid.New(),
		DateTime:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.ID, product.Type, product.ReceptionID, product.DateTime).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.AddProduct(context.Background(), product)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddProduct_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewProductRepo(mock)

	product := &models.Product{
		ID:          uuid.New(),
		Type:        "bag",
		ReceptionID: uuid.New(),
		DateTime:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.ID, product.Type, product.ReceptionID, product.DateTime).
		WillReturnError(errors.New("insert failed"))

	err = repo.AddProduct(context.Background(), product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось добавить продукт")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLastProduct_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewProductRepo(mock)

	pvzID := uuid.New()

	mock.ExpectExec("WITH last_product AS").
		WithArgs(pvzID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 row
	err = repo.DeleteLastProduct(context.Background(), pvzID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLastProduct_NoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewProductRepo(mock)

	pvzID := uuid.New()

	mock.ExpectExec("WITH last_product AS").
		WithArgs(pvzID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0)) // 0 rows

	err = repo.DeleteLastProduct(context.Background(), pvzID)
	assert.Error(t, err)
	assert.EqualError(t, err, "нет товаров для удаления")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLastProduct_ExecError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewProductRepo(mock)

	pvzID := uuid.New()

	mock.ExpectExec("WITH last_product AS").
		WithArgs(pvzID).
		WillReturnError(errors.New("delete failed"))

	err = repo.DeleteLastProduct(context.Background(), pvzID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось удалить последний товар")
	assert.NoError(t, mock.ExpectationsWereMet())
}
