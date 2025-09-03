package db

import (
	"fmt"
	"log"
	"os"
	// "log"
	// "os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
	initErr  error
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	MaxIdleConnection int
	MaxOpenConnection int
	ConnectionLifeTimeMinute int
	Logging bool
}

// Get singleton database connection
// func ConnectDb(cfg DBConfig) (*gorm.DB, error) {
// 	once.Do(func() {

// 		dsn := fmt.Sprintf(
// 			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
// 			cfg.DB.Host, cfg.DB.User, cfg.DB.Password,
// 			cfg.DB.DBName, cfg.DB.Port,
// 		)

// 		var logLevel logger.LogLevel
// 		if cfg.DB.Logging {
// 			logLevel = logger.Info // all logs
// 		} else {
// 			logLevel = logger.Silent
// 		}

// 		gormLogger := logger.New(
// 			log.New(os.Stdout, "\r\n", log.LstdFlags),
// 			logger.Config{
// 				SlowThreshold:             time.Second, // Log slower queries
// 				LogLevel:                  logLevel,    // Control log level
// 				IgnoreRecordNotFoundError: true,
// 				Colorful:                  true,
// 			},
// 		)

// 		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 			Logger: gormLogger,
// 		})
// 		if err != nil {
// 			initErr = fmt.Errorf("failed to connect to database: %w", err)
// 			return
// 		}

// 		sqlDB, err := db.DB()
// 		if err != nil {
// 			initErr = fmt.Errorf("failed to get underlying sql.DB: %w", err)
// 			return
// 		}

// 		sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConnection)
// 		sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConnection)
// 		sqlDB.SetConnMaxLifetime(time.Duration(cfg.DB.ConnectionLifeTimeMinute) * time.Minute)
// 		instance = db
// 	})
// 	fmt.Println("Database connected successfully")
// 	return instance, initErr
// }

func ConnectDb(cfg *DBConfig) (*gorm.DB, error) {
	once.Do(func() {
		var dsn string
		var dialector gorm.Dialector

		// Choose driver based on config
		switch cfg.Driver {
		case "mysql":
			dsn = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.User, cfg.Password,
				cfg.Host, cfg.Port,
				cfg.DBName,
			)
			dialector = mysql.Open(dsn)

		case "postgres":
			dsn = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				cfg.Host, cfg.User, cfg.Password,
				cfg.DBName, cfg.Port,
			)
			dialector = postgres.Open(dsn)

		case "mssql":
			dsn = fmt.Sprintf(
				"sqlserver://%s:%s@%s:%d?database=%s",
				cfg.User, cfg.Password,
				cfg.Host, cfg.Port,
				cfg.DBName,
			)
			dialector = sqlserver.Open(dsn)

		default:
			initErr = fmt.Errorf("unsupported database driver: %s", cfg.Driver)
			return
		}

		// Configure logger
		var logLevel logger.LogLevel
		if cfg.Logging {
			logLevel = logger.Info
		} else {
			logLevel = logger.Silent
		}

		gormLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		// Connect
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("failed to get underlying sql.DB: %w", err)
			return
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnectionLifeTimeMinute) * time.Minute)

		instance = db
	})
	fmt.Println("Your",cfg.Driver,"Database connected successfully")
	return instance, initErr
}

// CloseDb closes the singleton database connection gracefully
func CloseDb() error {
	if instance == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := instance.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	instance = nil
	return nil
}
