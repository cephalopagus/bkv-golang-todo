package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/cephalopagus/bkv-golang-todo/internal/core/config"
	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_pgx_pool "github.com/cephalopagus/bkv-golang-todo/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/middleware"
	core_http_server "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/cephalopagus/bkv-golang-todo/internal/features/tasks/repository/postgres"
	tasks_service "github.com/cephalopagus/bkv-golang-todo/internal/features/tasks/service"
	tasks_transport_http "github.com/cephalopagus/bkv-golang-todo/internal/features/tasks/transport/http"
	users_postrgres_repository "github.com/cephalopagus/bkv-golang-todo/internal/features/users/repository/postrgres"
	users_service "github.com/cephalopagus/bkv-golang-todo/internal/features/users/service"
	users_transport_http "github.com/cephalopagus/bkv-golang-todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postrgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransport := users_transport_http.NewUserHTTPHandler(usersService)

	//
	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	taskService := tasks_service.NewTasksService(tasksRepository)
	tasksTransport := tasks_transport_http.NewTaskHTTPHandler(taskService)

	logger.Debug("initializing http server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVer := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVer.RegisterRoutes(usersTransport.Routes()...)
	apiVer.RegisterRoutes(tasksTransport.Route()...)

	// apiVer2 := core_http_server.NewAPIVersionRouter(
	// 	core_http_server.ApiVersion2,
	// 	core_http_middleware.Dummy("api version 2 middleware"))

	// apiVer2.RegisterRoutes(usersTransport.Routes()...)

	httpServer.RegisterApiRouters(
		apiVer,
		// apiVer2,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}

}
