package external

import (
	"fmt"
	"time"

	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const (
	maxOpenConnNum  = 25
	maxIdleConnNum  = 25
	connMaxLifetime = 5 * time.Minute
)

type externalDB struct {
	client *gorm.DB
}

func (e *externalDB) DB() *gorm.DB {
	return e.client
}

func (e *externalDB) NewTxDB(tx *gorm.DB) domain.ExternalDBClient {
	return &externalDB{client: tx}
}

func MustNewExternalDB() domain.ExternalDBClient {
	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		envs.DBHost, envs.DBUser, envs.DBPassword, envs.DBName, envs.DBPort,
	)

	// Initialize GORM DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 pkgLogger.ZapLogger.GormLogger,
	})
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("failed to connect to database: %v", err)
	}

	if err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{postgres.Open(dsn)},
		Replicas: []gorm.Dialector{postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			envs.DBReplicaHost, envs.DBUser, envs.DBPassword, envs.DBName, envs.DBReplicaPort,
		))},
		Policy: dbresolver.RandomPolicy{},
	})); err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("DB resolver initialization: %v", err)
	}

	// Configure a connection pool
	sqlDB, err := db.DB()
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("failed to get database object: %v", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConnNum)
	sqlDB.SetMaxIdleConns(maxIdleConnNum)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	pkgLogger.ZapLogger.Logger.Info("PostgreSQL connection established successfully!")

	return &externalDB{client: db}
}
