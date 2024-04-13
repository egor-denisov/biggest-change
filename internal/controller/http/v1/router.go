package v1

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	_ "github.com/egor-denisov/biggest-change/docs"
	"github.com/egor-denisov/biggest-change/internal/controller/http/v1/jsonrpc"
	"github.com/egor-denisov/biggest-change/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger spec:
// @title       Stats Of Changing
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1 .
func NewRouter(handler *gin.Engine, l *slog.Logger, sc usecase.StatsOfChanging) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// JSON RPC
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterService(jsonrpc.NewStatsOfChangingService(l, sc), "JsonRpc")

	handler.POST("/", gin.WrapH(rpcServer))

	// Routers
	h := handler.Group("/api/v1")
	{
		newStatsOfChanging(h, l, sc)
	}
}
