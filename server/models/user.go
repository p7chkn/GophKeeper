package models

// User структура для хранения модели пользователя
type User struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}
