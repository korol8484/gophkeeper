package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"regexp"
	"testing"
	"time"
)

func TestSecretRepository_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	ctx := context.Background()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO secret (.+)`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg())

	repo := NewSecretRepository(db)

	repo.Add(ctx, nil, nil, 1, 1, time.Now())
}

func TestSecretRepository_GetAllByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	ctx := context.Background()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "meta_data", "context", "version", "added", "updated"})
	rows.AddRow(uuid.New(), "{}", "", 1, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, meta_data, context, version, added, updated FROM secret`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	repo := NewSecretRepository(db)

	repo.GetAllByUserID(ctx, 1)
}

func TestSecretRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	ctx := context.Background()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "meta_data", "context", "version", "added", "updated"})
	rows.AddRow(uuid.New(), "{}", "", 1, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, meta_data, context, version, added, updated FROM secret`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	repo := NewSecretRepository(db)

	repo.Get(ctx, 1, uuid.New())
}
