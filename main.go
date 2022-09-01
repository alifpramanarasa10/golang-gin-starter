package main

import (
	"context"
	"fmt"
	"gin-starter/app"
	"gin-starter/config"
	authBuilder "gin-starter/modules/auth/v1/builder"
	masterBuilder "gin-starter/modules/master/v1/builder"
	notificationBuilder "gin-starter/modules/notification/v1/builder"
	userBuilder "gin-starter/modules/user/v1/builder"
	pubsubSDK "gin-starter/sdk/pubsub"
	"gin-starter/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"gorm.io/gorm"
)

const (
	statusOK = "OK"
)

var environment string

// Health is a base struct for health check
type Health struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

var health *Health

// splash print plain text message to console
func splash() {
	fmt.Print(`
        .__                   __                 __                
   ____ |__| ____     _______/  |______ ________/  |_  ___________ 
  / ___\|  |/    \   /  ___/\   __\__  \\_  __ \   __\/ __ \_  __ \
 / /_/  >  |   |  \  \___ \  |  |  / __ \|  | \/|  | \  ___/|  | \/
 \___  /|__|___|  / /____  > |__| (____  /__|   |__|  \___  >__|   
/_____/         \/       \/            \/                 \/       
`)
}

func main() {
	health = &Health{}
	cfg, err := config.LoadConfig(".env")
	checkError(err)

	splash()

	environment = cfg.Env

	db, err := utils.NewPostgresGormDB(&cfg.Postgres)
	checkError(err)
	health.Database = statusOK

	redisPool := buildRedisPool(cfg)
	health.Redis = statusOK

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	tracer, closer, _ := NewJaegerTracer(cfg.AppName, fmt.Sprintf("%s:%s", cfg.Jaeger.Address, cfg.Jaeger.Port))

	defer func() {
		if err := closer.Close(); err != nil {
			log.Println("failed to close opentracing closer:", err)
		}
	}()

	opentracing.SetGlobalTracer(tracer)

	router.Use(OpenTracing())
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", Home)
	router.GET("/health-check", HealthGET)

	// Uncomment if you use AWS S3
	// awsSession := utils.InitAWSS3(*cfg)

	// Empty session just to handle error
	// Comment this line if you end up using aws s3
	awsSession := &session.Session{}

	BuildHandler(*cfg, router, db, redisPool, awsSession)

	// Uncomment if you use Google pub/sub
	//psClient := createPubSubClient(cfg.Google.ProjectID, cfg.Google.ServiceAccountFile)
	//psHandlers := registerPubsubHandlers(context.Background(), cfg, db, redisPool)
	//
	//_ = psClient.StartSubscriptions(psHandlers...)

	health.Status = statusOK
	if err := router.Run(fmt.Sprintf(":%s", cfg.Port.APP)); err != nil {
		panic(err)
	}
}

// BuildHandler is a function to build all handlers
func BuildHandler(cfg config.Config, router *gin.Engine, db *gorm.DB, redisPool *redis.Pool, awsSession *session.Session) {
	app.DefaultHTTPHandler(cfg, router)
	authBuilder.BuildAuthHandler(cfg, router, db, redisPool, awsSession)
	userBuilder.BuildUserHandler(cfg, router, db, redisPool, awsSession)
	notificationBuilder.BuildNotificationHandler(cfg, router, db, redisPool, awsSession)
	masterBuilder.BuildMasterHandler(cfg, router, db, redisPool, awsSession)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// CORSMiddleware is a function to add CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

// registerPubsubHandlers is a function to register pubsub handlers
func registerPubsubHandlers(
	ctx context.Context,
	cfg *config.Config,
	gconn *gorm.DB,
	redisPool *redis.Pool,
) []pubsubSDK.Subscriber {
	var handlers []pubsubSDK.Subscriber

	handlers = append(handlers, notificationBuilder.BuildSendEmailPubsubHandler(*cfg, gconn))
	return handlers
}

// buildRedisPool is a function to build redis pool
func buildRedisPool(cfg *config.Config) *redis.Pool {
	cachePool := utils.NewPool(cfg.Redis.Address, cfg.Redis.Password)

	ctx := context.Background()
	_, err := cachePool.GetContext(ctx)

	if err != nil {
		checkError(err)
	}

	log.Print("redis successfully connected!")
	return cachePool
}

// createPubSubClient is a function to create pubsub client
func createPubSubClient(projectID, googleSaFile string) *pubsubSDK.PubSub {
	return pubsubSDK.NewPubSub(projectID, &googleSaFile)
}

// NewJaegerTracer is a function to create a new Jaeger tracer
func NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHostPort,
		},

		ServiceName: serviceName,
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	fmt.Println("tracer:", tracer)

	return tracer, closer, err
}

// OpenTracing is a function to add OpenTracing middleware
func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		wireCtx, _ := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))

		serverSpan := opentracing.StartSpan(c.Request.URL.Path,
			ext.RPCServerOption(wireCtx))
		defer serverSpan.Finish()
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), serverSpan))
		c.Next()
	}
}

// HealthGET is a function to handle health check
func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, health)
}

// Home is a function to handle home page
func Home(c *gin.Context) {
	appVersion := os.Getenv("APP_VERSION")

	c.JSON(http.StatusOK, gin.H{
		"app_name":    "gin-starter-api",
		"environment": environment,
		"version":     appVersion,
		"status":      "running",
	})
}
