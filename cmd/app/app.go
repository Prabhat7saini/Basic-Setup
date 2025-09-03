package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	// "github.com/Prabhat7saini/TestingGo/shared/clients/db"
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/db"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/redis"

	// "github.com/Prabhat7saini/Basic-Setup/shared/clients/s3"
	"github.com/Prabhat7saini/Basic-Setup/shared/socket"

	// sharedRedis "github.com/Prabhat7saini/TestingGo/shared/clients/redis"

	// "github.com/Prabhat7saini/TestingGo/shared/clients/socket"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg     *config.Env
	router  *gin.Engine
	server  *http.Server
	db      *gorm.DB
	redis   redis.Client
	log     *zap.Logger
	// s3 s3.Client
	// socketHub *socket.Hub
	socketSrv *http.Server
}

func NewApp(cfg *config.Env, log *zap.Logger) *App {
	app := &App{
		cfg: cfg,
		log: log,
	}
	app.initialize()
	return app
}

func (a *App) initialize() {
	// Connect DB
	// dbConn, err := db.ConnectDb(Driver:"mysql", Host:"127.0.0.1", Port:"3306", User:"root", Password:"root", Database:"test", MaxIdleConns:5, MaxOpenConns:10, PoolSize:10, ConnMaxLifetime:30*time.Second, Log:log)
	// if err != nil {
	// 	a.log.Fatal("failed to connect to database", zap.Error(err))
	// }
	// a.db = dbConn
	dbConn, err := db.ConnectDb(&db.DBConfig{
    Driver: "postgres",
    Host: "localhost",
    Port: 5432,
    User: "prabhat",
    Password: "prabhat",
    DBName: "ChatAppGo",
    MaxIdleConnection: 10,
    MaxOpenConnection: 50,
    ConnectionLifeTimeMinute: 30,
    Logging: true,
})
	if err != nil {
		a.log.Fatal("failed to initialize database", zap.Error(err))
	}
	a.db = dbConn

	// Connect Redis
	redisConn, err := redis.InitRedis(a.cfg)
	if err != nil {
		a.log.Fatal("failed to initialize redis", zap.Error(err))
	}
	a.redis = redisConn

	// Init Router
	a.router = gin.Default()

	// routes.NewRoutes(a.router, a.handler)

	// validator.RegisterValidations()
	// Init HTTP server
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      a.router,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}

	// inside app.initialize()
	ioServer := socket.GetSocketIOServer()
	go ioServer.Serve()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", ioServer) // <--- socket.io endpoint
	a.socketSrv = &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	a.registerRoutes()
	

}

func (a *App) registerRoutes() {
	// Example route
	a.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

}

func (a *App) Run() {
	// Start Gin HTTP server
	go func() {
		a.log.Info("starting HTTP server", zap.String("addr", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Start Socket.IO server
	go func() {
		a.log.Info("starting WebSocket (Socket.IO) server", zap.String("addr", a.socketSrv.Addr))
		if err := a.socketSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("WebSocket server error", zap.Error(err))
		}
	}()

	// Wait for OS interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Signal(syscall.SIGTERM))
	<-quit
	a.log.Info("Received shutdown signal, shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("HTTP server forced to shutdown", zap.Error(err))
	} else {
		a.log.Info("HTTP server stopped gracefully")
	}

	// Shutdown Socket.IO server
	socket.CloseSocketIOServer()
	if err := a.socketSrv.Shutdown(ctx); err != nil {
		a.log.Error("Socket.IO server forced to shutdown", zap.Error(err))
	} else {
		a.log.Info("Socket.IO server stopped gracefully")
	}

	// Shutdown DB
	if err := db.CloseDb(); err != nil {
		a.log.Error("Database close error", zap.Error(err))
	} else {
		a.log.Info("Database connection closed gracefully")
	}

	// Shutdown Redis
	if err := redis.CloseRedis(); err != nil {
		a.log.Error("Redis close error", zap.Error(err))
	} else {
		a.log.Info("Redis connection closed gracefully")
	}
}
