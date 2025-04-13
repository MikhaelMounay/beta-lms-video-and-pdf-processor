package config

import (
	"context"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	SECRET_KEY     string
	IV             string
	VOD_KEY_ID     string
	VOD_KEY        string
	INSTANCE_NAME  string
	BASE_B2_FOLDER string
	ENDPOINT       string
	APP_KEY_ID     string
	APP_KEY        string
	BUCKET_ID      string
	BUCKET_NAME    string
	BUCKET_REGION  string
	DATABASE_URL   string
	DbPool         *pgxpool.Pool
)

func ConnectDB(log *logger.Logger) func() {
	var err error
	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	DbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		log.Fatal("Failed to connect to the database: %v", err)
	}

	return func() {
		DbPool.Close()
	}
}
