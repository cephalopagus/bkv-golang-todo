package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_postgres_pool "github.com/cephalopagus/bkv-golang-todo/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/middleware"
	core_http_server "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/server"
	users_postrgres_repository "github.com/cephalopagus/bkv-golang-todo/internal/features/users/repository/postrgres"
	users_service "github.com/cephalopagus/bkv-golang-todo/internal/features/users/service"
	users_transport_http "github.com/cephalopagus/bkv-golang-todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postrgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransport := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initializing http server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVer := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVer.RegisterRoutes(usersTransport.Routes()...)
	httpServer.RegisterApiRouters(apiVer)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}

}
