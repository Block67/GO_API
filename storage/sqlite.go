package storage

import (
	"fmt"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// SQLiteStorage gère la connexion à la base de données SQLite
type SQLiteStorage struct {
	Container *sqlstore.Container
}

// NewSQLiteStorage initialise une nouvelle connexion à la base de données SQLite
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath), dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database: %v", err)
	}
	return &SQLiteStorage{Container: container}, nil
}

// GetDevice récupère le store de l'appareil
func (s *SQLiteStorage) GetDevice() (*store.Device, error) {
	device := s.Container.NewDevice()
	return device, nil
}
