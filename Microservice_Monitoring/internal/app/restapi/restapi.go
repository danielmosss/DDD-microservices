package restapi

import (
	v1 "monitoring/internal/app/restapi/v1"
	v2 "monitoring/internal/app/restapi/v2"

	_ "monitoring/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartRestAPI() string {
	router := gin.Default()

	api := router.Group("/api")
	{
		v1Router := api.Group("/v1")
		{
			v1Router.GET("/status", v1.GetStatus)
		}

		v2Router := api.Group("/v2")
		{
			v2Router.GET("/status", v2.GetStatus)
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
