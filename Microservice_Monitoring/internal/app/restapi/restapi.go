package restapi

import (
	"monitoring/internal/app/restapi/frontend"
	v1 "monitoring/internal/app/restapi/v1"
	v2 "monitoring/internal/app/restapi/v2"
	"time"

	_ "monitoring/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartRestAPI() string {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-Requested-With", "Authorization", "Ngrok-Skip-Browser-Warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		v1Router := api.Group("/v1")
		{
			v1Router.GET("/status", v1.GetStatus)
			v1Router.GET("/metingen/:kunstwerkId", v1.GetMetingen)
			v1Router.GET("/metingen/:kunstwerkId/recent", v1.GetMetingenRecent)
			v1Router.GET("/afwijkingen/:kunstwerkId", v1.GetAfwijkingen)
			v1Router.POST("/meting", v1.PostMeting)
			v1Router.GET("/kunstwerken", v1.GetKunstwerken)
			v1Router.GET("/kunstwerken/:kunstwerkId/sensoren", v1.GetSensorenByKunstwerk)
			v1Router.GET("/kunstwerken/:kunstwerkId/dailyhealthupdate", v1.GetKunstwerkDHU)
		}

		v2Router := api.Group("/v2")
		{
			v2Router.GET("/status", v2.GetStatus)
		}

		frontendRouter := api.Group("/frontend")
		{
			frontendRouter.GET("/status", v2.GetStatus)
			frontendRouter.GET("/kunstwerken", v1.GetKunstwerken)
			frontendRouter.GET("/kunstwerken/:kunstwerkId/tree", frontend.GetKunstwerkTree)
			frontendRouter.POST("/kunstwerken/:kunstwerkId/sensoren/bulk-actueel", frontend.GetBulkSensorData)
			//frontendRouter.GET("/onderdelen/:id/sensoren", frontend.GetSensorenByOnderdeel)
			//frontendRouter.GET("/sensoren/:id/metingen?range=24h", frontend.GetMetingenForSensorid)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	// reroute api endpoint /docs to /swagger/index.html
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	err := router.Run(":8080")
	if err != nil {
		return "Fout bij starten van REST API: " + err.Error()
	}
	return "REST API gestart op :8080"
}
