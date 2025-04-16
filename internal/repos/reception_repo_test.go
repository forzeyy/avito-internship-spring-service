package repos_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateReception_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)

	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    uuid.New(),
		Status:   "in_progress",
		DateTime: time.Now(),
	}

	mock.ExpectExec("INSERT INTO receptions").
		WithArgs(reception.ID, reception.PVZID, reception.Status, reception.DateTime).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateReception(context.Background(), reception)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseLastReception_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)

	pvzID := uuid.New()
	now := time.Now()

	expected := models.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		Status:   "close",
		DateTime: now.Add(-time.Hour),
		ClosedAt: &now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "pvz_id", "status", "created_at", "closed_at",
	}).AddRow(expected.ID, expected.PVZID, expected.Status, expected.DateTime, expected.ClosedAt)

	mock.ExpectQuery("WITH last_reception AS").
		WithArgs(pvzID).
		WillReturnRows(rows)

	result, err := repo.CloseLastReception(context.Background(), pvzID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Status, result.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseLastReception_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)

	pvzID := uuid.New()

	mock.ExpectQuery("WITH last_reception AS").
		WithArgs(pvzID).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.CloseLastReception(context.Background(), pvzID)
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateReception_Failure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)

	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    uuid.New(),
		Status:   "in_progress",
		DateTime: time.Now(),
	}

	mock.ExpectExec("INSERT INTO receptions").
		WithArgs(reception.ID, reception.PVZID, reception.Status, reception.DateTime).
		WillReturnError(errors.New("insert failed"))

	err = repo.CreateReception(context.Background(), reception)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось создать приемку")
	assert.NoError(t, mock.ExpectationsWereMet())
}
