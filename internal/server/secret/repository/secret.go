package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"time"
)

const tableName = "secret"

type SecretRepository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

var secretAddSQL = `INSERT INTO %s (id, meta_data, context, user_id, version, added, updated) VALUES ($1, $2, $3, $4, $5, $6, $7);`

func (r *SecretRepository) Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID, version int, added time.Time) (*uuid.UUID, error) {
	insertId := uuid.New()
	metaDataJson, err := json.Marshal(metaData)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(
		ctx,
		fmt.Sprintf(secretAddSQL, tableName),
		insertId,
		metaDataJson,
		context,
		userId,
		version,
		added,
		added,
	).Err()
	if err != nil {

		return nil, err
	}

	return &insertId, nil
}

var secretGetByUserSQL = `SELECT id, meta_data, context, version, added, updated FROM %s WHERE user_id = $1;`

func (r *SecretRepository) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error) {
	var list []*model.Secret
	var metaDataJson []byte

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(secretGetByUserSQL, tableName), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var secret model.Secret
		err = rows.Scan(&secret.ID, &metaDataJson, &secret.Content, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(metaDataJson, &secret.MetaData)
		if err != nil {
			return nil, err
		}

		list = append(list, &secret)
	}

	return list, nil
}

var secretGetSQL = `SELECT id, meta_data, context, version, added, updated FROM %s WHERE user_id = $1 and id = $2;`

func (r *SecretRepository) Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error) {
	var secret model.Secret
	var metaDataJson []byte

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(secretGetSQL, tableName), userID, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&secret.ID, &metaDataJson, &secret.Content, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(metaDataJson, &secret.MetaData)
		if err != nil {
			return nil, err
		}
	}

	return &secret, nil
}
