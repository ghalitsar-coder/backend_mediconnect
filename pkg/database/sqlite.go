package database

import (
  "fmt"
  "log/slog"

  "github.com/glebarez/sqlite"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
)

func ConnectSQLite(dsn string) (*gorm.DB, error) {
  db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
  })
  if err != nil {
    return nil, fmt.Errorf("failed to connect to database: %w", err)
  }

  slog.Info("Successfully connected to SQLite using GORM")
  return db, nil
}
