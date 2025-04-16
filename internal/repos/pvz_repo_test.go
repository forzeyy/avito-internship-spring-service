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

func TestCreatePVZ_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewPVZRepo(mock)

	pvz := &models.PVZ{
		ID:      uuid.New(),
		City:    "Москва",
		RegDate: time.Now(),
	}

	mock.ExpectExec("INSERT INTO pvzs").
		WithArgs(pvz.ID, pvz.City, pvz.RegDate).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreatePVZ(context.Background(), pvz)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePVZ_WrongCity(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewPVZRepo(mock)

	pvz := &models.PVZ{
		ID:      uuid.New(),
		City:    "Курск",
		RegDate: time.Now(),
	}

	mock.ExpectExec("INSERT INTO pvzs").
		WithArgs(pvz.ID, pvz.City, pvz.RegDate).
		WillReturnError(errors.New("insert failed"))

	err = repo.CreatePVZ(context.Background(), pvz)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось создать ПВЗ")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPVZs_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewPVZRepo(mock)

	start := time.Now().Add(-30 * 24 * time.Hour)
	end := time.Now()

	expected := models.PVZ{
		ID:      uuid.New(),
		City:    "Москва",
		RegDate: time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "city", "reg_date"}).
		AddRow(expected.ID, expected.City, expected.RegDate)

	mock.ExpectQuery("SELECT id, city, reg_date FROM pvzs").
		WithArgs(&start, &end, 10, 0).
		WillReturnRows(rows)

	result, err := repo.GetPVZs(context.Background(), &start, &end, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expected.City, result[0].City)
	assert.NoError(t, mock.ExpectationsWereMet())
}
