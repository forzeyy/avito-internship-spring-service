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

// CreateReception
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

// CloseLastReception
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

func TestGetLastOpenReception_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)

	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	rows := pgxmock.NewRows([]string{"id", "pvz_id", "status", "created_at", "closed_at"}).
		AddRow(receptionID, pvzID, "in_progress", now, nil)

	mock.ExpectQuery("SELECT id, pvz_id, status, created_at, closed_at").
		WithArgs(pvzID).
		WillReturnRows(rows)

	rec, err := repo.GetLastOpenReception(ctx, pvzID)

	assert.NoError(t, err)
	assert.NotNil(t, rec)
	assert.Equal(t, receptionID, rec.ID)
	assert.Equal(t, "in_progress", rec.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastOpenReception_NoRows(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)
	pvzID := uuid.New()

	mock.ExpectQuery("SELECT id, pvz_id, status, created_at, closed_at").
		WithArgs(pvzID).
		WillReturnError(pgx.ErrNoRows)

	rec, err := repo.GetLastOpenReception(ctx, pvzID)

	assert.NoError(t, err)
	assert.Nil(t, rec)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastOpenReception_QueryError(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewReceiptRepo(mock)
	pvzID := uuid.New()

	mock.ExpectQuery("SELECT id, pvz_id, status, created_at, closed_at").
		WithArgs(pvzID).
		WillReturnError(errors.New("unexpected db error"))

	rec, err := repo.GetLastOpenReception(ctx, pvzID)

	assert.Error(t, err)
	assert.Nil(t, rec)
	assert.Contains(t, err.Error(), "не удалось получить последнюю открытую приемку")
	assert.NoError(t, mock.ExpectationsWereMet())
}
