package builder

import (
	"gin-starter/app"
	"gin-starter/config"
	"gin-starter/modules/activities/v1/repository"
	"gin-starter/modules/activities/v1/service"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// BuildActivitiesHandler builds activities handler
// starting from handler down to repository or tool.
func BuildActivitiesHandler(cfg config.Config, router *gin.Engine, db *gorm.DB, redisPool *redis.Pool, awsSession *session.Session) {
	// Repository
	ar := repository.NewActivitiesRepository(db)

	// Service
	ac := service.NewActivitiesCreator(cfg, ar)
	af := service.NewActivitiesFinder(cfg, ar)
	ad := service.NewActivitiesDeleter(cfg, ar)
	au := service.NewActivitiesUpdater(cfg, ar)

	// Handler
	app.ActivitiesFinderHTTPHandler(cfg, router, af)
	app.ActivitiesCreatorHTTPHandler(cfg, router, ac)
	app.ActivitiesDeleterHTTPHandler(cfg, router, ad)
	app.ActivitiesUpdaterHTTPHandler(cfg, router, au)
}
