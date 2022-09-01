package builder

import (
	"gin-starter/app"
	"gin-starter/config"
	authRepo "gin-starter/modules/auth/v1/repository"
	auth "gin-starter/modules/auth/v1/service"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// BuildAuthHandler build auth handlers
// starting from handler down to repository or tool.
func BuildAuthHandler(cfg config.Config, router *gin.Engine, db *gorm.DB, redisPool *redis.Pool, awsSession *session.Session) {
	// Repository
	ar := authRepo.NewAuthRepository(db)

	uc := auth.NewAuthService(cfg, ar)

	app.AuthHTTPHandler(cfg, router, uc)
}
