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

func (r *SecretRepository) Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID, version int, added time.Time) (*uuid.UUID, error) {
	var id *uuid.UUID

	insertId := uuid.New()
	metaDataJson, err := json.Marshal(metaData)

	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`INSERT INTO %s (id, meta_data, context, user_id, version, added, updated) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id;`, tableName),
		insertId,
		metaDataJson,
		context,
		userId,
		version,
		added,
		added,
	).Scan(&id)
	if err != nil {

		return nil, err
	}

	return id, nil
}

func (r *SecretRepository) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error) {
	var list []*model.Secret
	var secret model.Secret
	var metaDataJson []byte

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id, meta_data, context, version, added, updated FROM %s WHERE user_id = $1;`, tableName), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&secret.ID, &metaDataJson, &secret.Content, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			println(err.Error())
			continue
		}
		err = json.Unmarshal(metaDataJson, &secret.MetaData)
		if err != nil {
			println(err.Error())
			continue
		}

		list = append(list, &secret)
	}

	return list, nil
}

func (r *SecretRepository) Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error) {
	var secret model.Secret
	var metaDataJson []byte

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT id, meta_data, context, version, added, updated FROM %s WHERE user_id = $1 and id = $2;`, tableName), userID, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&secret.ID, &metaDataJson, &secret.Content, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			println(err.Error())
			continue
		}
		err = json.Unmarshal(metaDataJson, &secret.MetaData)
		if err != nil {
			println(err.Error())
			continue
		}
	}

	return &secret, nil
}
