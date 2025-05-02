package storage

import (
	"database/sql"
	"whatsapp-mougni-api-go/models"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(dbPath string) (*UserStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Créer la table users
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            number TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            api_key TEXT NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, err
	}

	return &UserStore{db: db}, nil
}

func (s *UserStore) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO users (name, number, email, password, api_key) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Number, user.Email, hashedPassword, user.APIKey,
	)
	return err
}

func (s *UserStore) GetUserByAPIKey(apiKey string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, name, number, email, api_key, created_at FROM users WHERE api_key = ?",
		apiKey,
	).Scan(&user.ID, &user.Name, &user.Number, &user.Email, &user.APIKey, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
