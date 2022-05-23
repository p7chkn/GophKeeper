// Package database пакет для работы с базами данных
package database

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	customerrors "new_diplom/errors"
	"new_diplom/models"
)

// NewPostgresDataBase функция по созданию нового объекта для работы с базой данных
func NewPostgresDataBase(conn *sqlx.DB) *PostgresDataBase {
	return &PostgresDataBase{
		conn: conn,
	}
}

// PostgresDataBase структура данных для подключения к базе данных
type PostgresDataBase struct {
	conn *sqlx.DB
}

// CreateUser функция для создание пользователя по модели пользователя
func (pg *PostgresDataBase) CreateUser(ctx context.Context, user models.User) error {

	_, err := pg.conn.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf', 8)))",
		user.Login, user.Password)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return customerrors.NewCustomError(err, "user already exists")
		}
	}
	return err
}

// CheckUserPassword функция проверки соответствия пароля, возвращает так же id пользователя
func (pg *PostgresDataBase) CheckUserPassword(ctx context.Context, user models.User) (string, error) {
	var result string
	query, err := pg.conn.QueryxContext(ctx, `SELECT id FROM users WHERE login = $1 
                       AND password = crypt($2, password) AND deleted_at is NULL FETCH FIRST ROW ONLY;`,
		user.Login, user.Password)
	if err != nil {
		return result, customerrors.NewCustomError(err, "can't find user")
	}
	query.Next()
	err = query.Scan(&result)
	if err != nil {
		return result, customerrors.NewCustomError(err, "error with getting user_id")
	}
	return result, err
}

// DeleteUser функция удаления пользователя
func (pg *PostgresDataBase) DeleteUser(ctx context.Context, userID string) error {
	_, err := pg.conn.ExecContext(ctx, `UPDATE users SET deleted_at = current_timestamp WHERE id = $1`, userID)
	if err != nil {
		return customerrors.NewCustomError(err, "error with deleting user")
	}
	return nil
}

// AddSecret функция дл добавления секрета пользователя
func (pg *PostgresDataBase) AddSecret(ctx context.Context, secret models.RawSecretData) error {
	_, err := pg.conn.ExecContext(ctx,
		"INSERT INTO secrets (user_id, secret_data) VALUES ($1, $2)",
		secret.UserID, secret.Data)
	return err
}

// GetSecrets функция для получение всех секретов пользователя
func (pg *PostgresDataBase) GetSecrets(ctx context.Context, userID string) ([]models.RawSecretData, error) {
	rows, err := pg.conn.QueryxContext(ctx, "SELECT id, user_id,secret_data FROM secrets WHERE user_id=$1 AND deleted_at IS NULL", userID)
	if err != nil {
		return nil, err
	}
	var result []models.RawSecretData
	for rows.Next() {
		m := models.RawSecretData{}
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, err
}

// DeleteSecret функция для удаления секрета пользователя
func (pg *PostgresDataBase) DeleteSecret(ctx context.Context, secretID string, userID string) error {
	_, err := pg.conn.ExecContext(ctx, `UPDATE secrets SET deleted_at = current_timestamp WHERE id = $1 AND user_id = $2`,
		secretID, userID)
	return err
}
