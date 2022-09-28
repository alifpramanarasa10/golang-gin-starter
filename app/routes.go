package app

import (
	"gin-starter/common/interfaces"
	"gin-starter/config"
	"gin-starter/middleware"
	activitieshandlerv1 "gin-starter/modules/activities/v1/handler"
	activitiesservicev1 "gin-starter/modules/activities/v1/service"
	authhandlerv1 "gin-starter/modules/auth/v1/handler"
	authservicev1 "gin-starter/modules/auth/v1/service"
	masterhandlerv1 "gin-starter/modules/master/v1/handler"
	masterservicev1 "gin-starter/modules/master/v1/service"
	notificationhandlerv1 "gin-starter/modules/notification/v1/handler"
	notificationservicev1 "gin-starter/modules/notification/v1/service"
	userhandlerv1 "gin-starter/modules/user/v1/handler"
	userservicev1 "gin-starter/modules/user/v1/service"
	"gin-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeprecatedAPI is a handler for deprecated APIs
func DeprecatedAPI(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.ErrorAPIResponse(http.StatusForbidden, "this version of api is deprecated. please use another version."))
	c.Abort()
}

// DefaultHTTPHandler is a handler for default APIs
func DefaultHTTPHandler(cfg config.Config, router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.ErrorAPIResponse(http.StatusNotFound, "invalid route"))
		c.Abort()
	})
}

// AuthHTTPHandler is a handler for auth APIs
func AuthHTTPHandler(cfg config.Config, router *gin.Engine, auc authservicev1.AuthUseCase) {
	hnd := authhandlerv1.NewAuthHandler(auc)
	v1 := router.Group("/v1")
	{
		v1.POST("/user/login", hnd.Login)
		v1.POST("/cms/login", hnd.LoginCMS)
	}
}

// NotificationFinderHTTPHandler is a handler for notification APIs
func NotificationFinderHTTPHandler(cfg config.Config, router *gin.Engine, cf notificationservicev1.NotificationFinderUseCase, nu notificationservicev1.NotificationUpdaterUseCase) {
	hnd := notificationhandlerv1.NewNotificationFinderHandler(cf, nu)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	{
		v1.GET("/user/notifications", hnd.GetNotification)
		v1.GET("/user/notification/count", hnd.CountUnreadNotifications)
	}
}

// NotificationCreatorHTTPHandler is a handler for notification APIs
func NotificationCreatorHTTPHandler(cfg config.Config, router *gin.Engine, cf notificationservicev1.NotificationCreatorUseCase) {
	hnd := notificationhandlerv1.NewNotificationCreatorHandler(cf)
	v1 := router.Group("/v1")
	{
		v1.POST("/cms/notification", hnd.CreateNotification)
	}
}

// NotificationUpdaterHTTPHandler is a handler for notification APIs
func NotificationUpdaterHTTPHandler(cfg config.Config, router *gin.Engine, cf notificationservicev1.NotificationUpdaterUseCase) {
	hnd := notificationhandlerv1.NewNotificationUpdaterHandler(cf)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	{
		v1.PUT("/user/notification/set", hnd.RegisterUnregisterPlayerID)
		v1.PUT("/user/notification/read", hnd.UpdateReadNotification)
	}
}

// MasterFinderHTTPHandler is a handler for master APIs
func MasterFinderHTTPHandler(cfg config.Config, router *gin.Engine, mf masterservicev1.MasterFinderUseCase) {
	hnd := masterhandlerv1.NewMasterFinderHandler(mf)
	v1 := router.Group("/v1")
	{
		v1.GET("/provinces", hnd.GetProvinces)
		v1.GET("/regencies/:province_id", hnd.GetRegenciesByProvinceID)
		v1.GET("/districts/:regency_id", hnd.GetDistrictsByRegencyID)
		v1.GET("/villages/:district_id", hnd.GetVillagesByDistrictID)
	}
}

// MasterCreatorHTTPHandler is a handler for master APIs
func MasterCreatorHTTPHandler(cfg config.Config, router *gin.Engine, mc masterservicev1.MasterCreatorUseCase, cloudStorage interfaces.CloudStorageUseCase) {
	_ = masterhandlerv1.NewMasterCreatorHandler(mc, cloudStorage)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
}

// UserFinderHTTPHandler is a handler for user APIs
func UserFinderHTTPHandler(cfg config.Config, router *gin.Engine, cf userservicev1.UserFinderUseCase) {
	hnd := userhandlerv1.NewUserFinderHandler(cf)
	v1 := router.Group("/v1")
	{
		v1.GET("/user/forgot-password/profile/:token", hnd.GetUserByForgotPasswordToken)
	}

	v1.Use(middleware.Auth(cfg))
	{
		v1.GET("/user/profile", hnd.GetUserProfile)
	}

	v1.Use(middleware.Admin(cfg))
	{
		v1.GET("/cms/profile", hnd.GetAdminProfile)
		v1.GET("/cms/admin/list", hnd.GetAdminUsers)
		v1.GET("/cms/admin/detail/:id", hnd.GetAdminUserByID)
		v1.GET("/cms/user/list", hnd.GetUsers)
		v1.GET("/cms/user/detail/:id", hnd.GetUserByID)
		v1.GET("/cms/roles", hnd.GetRoles)
		v1.GET("/cms/permission", hnd.GetPermissions)
		v1.GET("/cms/user/permission", hnd.GetUserPermissions)
	}
}

// UserCreatorHTTPHandler is a handler for user APIs
func UserCreatorHTTPHandler(cfg config.Config, router *gin.Engine, uc userservicev1.UserCreatorUseCase, uf userservicev1.UserFinderUseCase, cloudStorage interfaces.CloudStorageUseCase) {
	hnd := userhandlerv1.NewUserCreatorHandler(uc, cloudStorage)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
	{
		v1.POST("/cms/user", hnd.CreateUser)
		v1.POST("/cms/admin/user", hnd.CreateAdmin)
		v1.POST("/cms/permission", hnd.CreatePermission)
		v1.POST("/cms/role", hnd.CreateRole)
	}
}

// UserUpdaterHTTPHandler is a handler for user APIs
func UserUpdaterHTTPHandler(cfg config.Config, router *gin.Engine, uu userservicev1.UserUpdaterUseCase, uf userservicev1.UserFinderUseCase, cloudStorage interfaces.CloudStorageUseCase) {
	hnd := userhandlerv1.NewUserUpdaterHandler(uu, uf, cloudStorage)
	v1 := router.Group("/v1")
	{
		v1.PUT("/user/forgot-password/request", hnd.ForgotPasswordRequest)
		v1.PUT("/user/forgot-password", hnd.ForgotPassword)
	}

	v1.Use(middleware.Auth(cfg))
	{
		v1.PUT("/user/profile", hnd.UpdateUser)
		v1.PUT("/user/password", hnd.ChangePassword)
		v1.PUT("/verify/otp", hnd.VerifyOTP)
		v1.PUT("/resend/otp", hnd.ResendOTP)
	}

	v1.Use(middleware.Admin(cfg))
	{
		v1.PUT("/cms/admin/:id", hnd.UpdateAdmin)
		v1.PUT("/cms/user/activate/:id", hnd.ActivateDeactivateUser)
		v1.PUT("/cms/role/:id", hnd.UpdateRole)
		v1.PUT("/cms/permission/:id", hnd.UpdatePermission)
	}
}

// UserDeleterHTTPHandler is a handler for user APIs
func UserDeleterHTTPHandler(cfg config.Config, router *gin.Engine, ud userservicev1.UserDeleterUseCase, cloudStorage interfaces.CloudStorageUseCase) {
	hnd := userhandlerv1.NewUserDeleterHandler(ud, cloudStorage)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
	{
		v1.DELETE("/cms/admin/:id", hnd.DeleteAdmin)
		v1.DELETE("/cms/role/:id", hnd.DeleteRole)
	}
}

// ActivitiesFinderHTTPHandler is a handler for activities APIs
func ActivitiesFinderHTTPHandler(cfg config.Config, router *gin.Engine, af activitiesservicev1.ActivitiesFinderUseCase) {
	hnd := activitieshandlerv1.NewActivitiesFinderHandler(af)
	v1 := router.Group("/v1")
	{
		v1.GET("/activities", hnd.GetActivities)
		v1.GET("/activities/:id", hnd.GetActivityByID)
	}
}

// ActivitiesCreatorHTTPHandler is a handler for activities APIs
func ActivitiesCreatorHTTPHandler(cfg config.Config, router *gin.Engine, ac activitiesservicev1.ActivitiesCreatorUseCase) {
	hnd := activitieshandlerv1.NewActivitiesCreatorHandler(ac)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
	{
		v1.POST("/activities", hnd.CreateActivities)
	}
}

// ActivitiesUpdaterHTTPHandler is a handler for activities APIs
func ActivitiesUpdaterHTTPHandler(cfg config.Config, router *gin.Engine, au activitiesservicev1.ActivitiesUpdaterUseCase) {
	hnd := activitieshandlerv1.NewActivitiesUpdaterHandler(au)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
	{
		v1.PUT("/activities/:id", hnd.UpdateActivities)
	}
}

// ActivitiesDeleterHTTPHandler is a handler for activities APIs
func ActivitiesDeleterHTTPHandler(cfg config.Config, router *gin.Engine, ad activitiesservicev1.ActivitiesDeleterUseCase) {
	hnd := activitieshandlerv1.NewActivitiesDeleterHandler(ad)
	v1 := router.Group("/v1")

	v1.Use(middleware.Auth(cfg))
	v1.Use(middleware.Admin(cfg))
	{
		v1.DELETE("/activities/:id", hnd.DeleteActivities)
	}
}
