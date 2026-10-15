package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
	middleware        []core_http_middleware.Middleware
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
	middleware ...core_http_middleware.Middleware,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
		middleware:        middleware,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/statistics",
			Middleware: h.middleware,
			Handler:    h.GetStatistics,
		},
	}
}
