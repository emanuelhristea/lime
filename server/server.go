package server

import (
	"html/template"
	"net/http"

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
	r := gin.Default()
	r.Use(sessions.Sessions(config.Config().GetString("cookie_name"), sessions.NewCookieStore([]byte(config.Config().GetString("cookie_secret")))))
	r.Static("/assets", "./server/web/assets")
	r.SetFuncMap(template.FuncMap{
		"formatAsPrice": formatAsPrice,
		"formatAsCheck": formatAsCheck,
		"formatAsDate":  formatAsDate,
		"columnStatus":  columnStatus,
		"bytesToString": keyBytesToString,
	})
	r.LoadHTMLGlob("./server/web/templates/*.html")

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	api := r.Group("/api")
	api.POST("/key", controllers.CreateKey)
	api.GET("/key/:customer_id", controllers.GetKey)
	api.PATCH("/key/:customer_id", controllers.UpdateKey)
	api.POST("/verify", controllers.VerifyKey)
	api.Use(middleware.AuthRequired)
	{
		api.POST("/tariff", controllers.CreateTariff)
		api.DELETE("/tariff/:id", controllers.DeleteTariff)

		api.POST("/customer", controllers.CreateCustomer)
		api.PATCH("/customer/:id", controllers.UpdateCustomer)
		api.DELETE("/customer/:id", controllers.DeleteCustomer)

		api.POST("/subscriptions", controllers.CreateCustomer)
		api.PATCH("/customer/:id", controllers.UpdateCustomer)
		api.DELETE("/customer/:id", controllers.DeleteCustomer)
	}

	admin := r.Group("/admin")
	admin.GET("/", controllers.MainHandler)
	admin.POST("/login", middleware.Login)
	admin.POST("/logout", middleware.Logout)
	admin.Use(middleware.AuthRequired)
	{
		admin.GET("/license/:id", controllers.DownloadLicense)

		admin.GET("/subscription/:id/*action", controllers.CustomerSubscriptionList)
		admin.GET("/tariffs/*action", controllers.TariffsList)
		admin.GET("/customers/*action", controllers.MainHandler)
	}

	return r
}
