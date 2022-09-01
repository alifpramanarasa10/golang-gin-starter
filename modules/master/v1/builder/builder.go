package builder

import (
	"gin-starter/app"
	"gin-starter/config"
	"gin-starter/modules/master/v1/repository"
	"gin-starter/modules/master/v1/service"
	"gin-starter/sdk/gcs"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// BuildMasterHandler builds master handler
// starting from handler down to repository or tool.
func BuildMasterHandler(cfg config.Config, router *gin.Engine, db *gorm.DB, redisPool *redis.Pool, awsSession *session.Session) {
	// Repository
	pr := repository.NewProvinceRepository(db)
	dr := repository.NewDistrictRepository(db)
	vr := repository.NewVillageRepository(db)
	rr := repository.NewRegencyRepository(db)
	cloudStorage := gcs.NewGoogleCloudStorage(cfg)
	// cloudStorage := aws.NewS3Bucket(cfg, awsSession)

	// Service
	mc := service.NewMasterCreator(cfg, cloudStorage)
	mf := service.NewMasterFinder(cfg, pr, rr, dr, vr)

	// Handler
	app.MasterFinderHTTPHandler(cfg, router, mf)
	app.MasterCreatorHTTPHandler(cfg, router, mc, cloudStorage)
}
