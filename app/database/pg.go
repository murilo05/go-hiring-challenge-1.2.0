package database

import (
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PG struct {
	*gorm.DB
	logger *zap.SugaredLogger
}

func New(user, password, dbname, port string, logger *zap.SugaredLogger) (pg *PG, close func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect database: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Failed to get database connection: ", err)
	}

	return &PG{
		db,
		logger,
	}, sqlDB.Close
}
