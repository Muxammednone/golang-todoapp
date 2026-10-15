package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Muxammednone/golang-todoapp/internal/core/config"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/Muxammednone/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/server"
	auth_client_grpc "github.com/Muxammednone/golang-todoapp/internal/features/auth/client/grpc"
	auth_token "github.com/Muxammednone/golang-todoapp/internal/features/auth/token"
	auth_transport_http "github.com/Muxammednone/golang-todoapp/internal/features/auth/transport/http"
	statistics_postgres_repository "github.com/Muxammednone/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Muxammednone/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/Muxammednone/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Muxammednone/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Muxammednone/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Muxammednone/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Muxammednone/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Muxammednone/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Muxammednone/golang-todoapp/internal/features/users/transport/http"
	web_fs_repository "github.com/Muxammednone/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/Muxammednone/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/Muxammednone/golang-todoapp/internal/features/web/transport/http"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/Muxammednone/golang-todoapp/docs"
)

// @title 			Golang Todo API
// @version 		1.0
// @description 	Todo Application REST-API schema
// @BasePath 		/api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter Bearer <access_token> returned by POST /auth/login.

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
		fmt.Println("failed to init application logger")
		os.Exit(1)
	}

	defer logger.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		logger.Error("JWT_SECRET must be at least 32 bytes")
		return
	}
	verifier := auth_token.NewVerifier([]byte(jwtSecret))
	authMiddleware := core_http_middleware.Auth(verifier)

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, authMiddleware)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTaskRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService, authMiddleware, core_http_middleware.RequireProfile(usersService))

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statistivsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statistivsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService, authMiddleware, core_http_middleware.RequireProfile(usersService))

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	ssoConfig := auth_client_grpc.NewConfigMust()
	ssoConn, err := grpc.NewClient(
		ssoConfig.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to create SSO client", zap.Error(err))
		return
	}
	defer ssoConn.Close()
	authClient := auth_client_grpc.New(ssoConn, ssoConfig.Timeout)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authClient, authMiddleware)
	logger.Debug("Initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)

	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
