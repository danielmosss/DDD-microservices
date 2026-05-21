package restapi

import (
	v1 "monitoring/internal/app/restapi/v1"
	v2 "monitoring/internal/app/restapi/v2"

	"github.com/gin-gonic/gin"
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

	err := router.Run(":8080")
	if err != nil {
		return "Fout bij starten van REST API: " + err.Error()
	}
	return "REST API gestart op :8080"
}
