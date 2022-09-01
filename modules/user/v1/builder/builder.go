package builder

import (
	"gin-starter/app"
	"gin-starter/config"
	notificationRepo "gin-starter/modules/notification/v1/repository"
	notification "gin-starter/modules/notification/v1/service"
	userRepo "gin-starter/modules/user/v1/repository"
	"gin-starter/modules/user/v1/service"
	"gin-starter/sdk/gcs"
	"gin-starter/utils"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// BuildUserHandler builds user handler
// starting from handler down to repository or tool.
func BuildUserHandler(cfg config.Config, router *gin.Engine, db *gorm.DB, redisPool *redis.Pool, awsSession *session.Session) {
	// Cache
	cache := utils.NewClient(redisPool)

	// Repository
	ur := userRepo.NewUserRepository(db)
	rr := userRepo.NewRoleRepository(db, cache)
	urr := userRepo.NewUserRoleRepository(db, cache)
	pr := userRepo.NewPermissionRepository(db, cache)
	nr := notificationRepo.NewNotificationRepository(db)
	
	// Cloud Storage
	cloudStorage := gcs.NewGoogleCloudStorage(cfg)
	// cloudStorage := aws.NewS3Bucket(cfg, awsSession)

	// Service
	nc := notification.NewNotificationCreator(cfg, nr)
	uc := service.NewUserCreator(cfg, ur, urr, rr, pr, nc, cloudStorage)
	uf := service.NewUserFinder(cfg, ur, urr, rr, pr)
	uu := service.NewUserUpdater(cfg, ur, urr, rr, pr)
	ud := service.NewUserDeleter(cfg, ur, rr)

	// Handler
	app.UserFinderHTTPHandler(cfg, router, uf)
	app.UserCreatorHTTPHandler(cfg, router, uc, uf, cloudStorage)
	app.UserUpdaterHTTPHandler(cfg, router, uu, uf, cloudStorage)
	app.UserDeleterHTTPHandler(cfg, router, ud, cloudStorage)
}
