package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/emanuelhristea/lime/config"
	"github.com/emanuelhristea/lime/server/controllers"
	"github.com/emanuelhristea/lime/server/middleware"
	"github.com/emanuelhristea/lime/server/seed"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Start is a ...
func Start(args []string) {
	cfg := config.Config()
	gin.SetMode(cfg.GetString("mode"))
	r := setupRouter()

	for _, arg := range args {
		switch arg {
		case "seed":
			seed.Load(config.DB)
		}
	}

	err := r.Run(cfg.GetString("port"))
	if err != nil {
		panic(err)
	}
}

func setupRouter() *gin.Engine {
	cfg := config.Config()
	webPath := cfg.GetString("web_path")
	r := gin.Default()
	r.Use(sessions.Sessions(cfg.GetString("cookie_name"), sessions.NewCookieStore([]byte(cfg.GetString("cookie_secret")))))
	r.Static("/assets", webPath+"/assets")
	r.SetFuncMap(template.FuncMap{
		"formatAsPrice":         formatAsPrice,
		"formatAsDateTime":      formatAsDateTime,
		"formatAsDateTimeLocal": formatAsDateTimeLocal,
		"columnStatus":          columnStatus,
		"bytesToString":         keyBytesToString,
	})
	r.LoadHTMLGlob(webPath + "/templates/*.html")

	// Create limiter struct
	limiterOptions := &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
		ExpireJobInterval:    time.Hour * 2,
	}

	limiterSlow := tollbooth_gin.LimitHandler(tollbooth.NewLimiter(5, limiterOptions))
	limiterMedium := tollbooth_gin.LimitHandler(tollbooth.NewLimiter(10, limiterOptions))
	limiterFast := tollbooth_gin.LimitHandler(tollbooth.NewLimiter(30, limiterOptions))

	r.GET("/ping", limiterSlow, func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	api := r.Group("/api")
	api.POST("/key", limiterFast, controllers.CreateKey)
	api.DELETE("/key", limiterFast, controllers.ReleaseKey)
	api.GET("/subscriptions", limiterFast, controllers.GetUserSubscriptions)
	api.PATCH("/key/:customer_id", limiterFast, controllers.UpdateLicense)
	api.POST("/verify", limiterFast, controllers.VerifyKey)
	api.Use(middleware.AuthRequired)
	{
		api.GET("/tariffs", limiterFast, controllers.GetTariffList)
		api.GET("/tariff/:id", limiterFast, controllers.GetTariff)
		api.POST("/tariff", limiterMedium, controllers.CreateTariff)
		api.PATCH("/tariff/:id", limiterMedium, controllers.UpdateTariff)
		api.DELETE("/tariff/:id", limiterSlow, controllers.DeleteTariff)

		api.GET("/customers", limiterFast, controllers.GetCustomerList)
		api.GET("/customer/:id", limiterFast, controllers.GetCustomer)
		api.POST("/customer", limiterMedium, controllers.CreateCustomer)
		api.PATCH("/customer/:id", limiterMedium, controllers.UpdateCustomer)
		api.DELETE("/customer/:id", limiterSlow, controllers.DeleteCustomer)

		api.GET("/subscriptions/:customerId", limiterFast, controllers.GetSubscriptionList)
		api.GET("/subscription/:id", limiterFast, controllers.GetSubscription)
		api.PATCH("/subscription/:sid/renew", limiterMedium, controllers.ReNewSubscription)
		api.POST("/customer/:id/subscription", limiterMedium, controllers.CreateSubscription)
		api.PATCH("/customer/:id/subscription/:sid", limiterMedium, controllers.UpdateSubscription)
		api.DELETE("/customer/:id/subscription/:sid", limiterSlow, controllers.DeleteSubscription)

		api.GET("/licenses/:subscripionId", limiterFast, controllers.GetLicensesList)
		api.GET("/license/:id", limiterFast, controllers.GetLicense)
		api.POST("/licenses/:subscripionId", limiterFast, controllers.CreateLicense)
		api.PATCH("/license/:id", limiterMedium, controllers.UpdateLicense)
		api.DELETE("/license/:id", limiterSlow, controllers.DeleteLicense)
	}

	admin := r.Group("/admin")
	admin.GET("/", controllers.MainHandler)
	admin.POST("/login", middleware.Login)
	admin.POST("/logout", middleware.Logout)
	admin.Use(middleware.AuthRequired)
	{
		admin.GET("/license/:id", limiterMedium, controllers.DownloadLicense)

		admin.GET("/customer/:id", controllers.CustomerRowHandler)
		admin.GET("/customer/:id/subscription/*action", controllers.CustomerSubscriptionAction)
		admin.GET("/customer/:id/subscriptions/*action", controllers.CustomerSubscriptionList)

		admin.GET("/customer/:id/sub/:sid/subscription/*action", controllers.CustomerSubscriptionAction)
		admin.GET("/customer/:id/sub/:sid/license/*action", controllers.CustomerSubscriptionLicenseAction)

		admin.GET("/tariffs/*action", controllers.TariffsList)
		admin.GET("/tariff/:id/edit/*action", controllers.TariffAction)

		admin.GET("/customers/*action", controllers.MainHandler)
		admin.GET("/customer/:id/edit/*action", controllers.MainHandler)
	}

	return r
}
