package api

import (
	"github.com/gin-gonic/gin"
)

func Route(e *gin.Engine) {
	g := e.Group("/api")

	collection := g.Group("/collection")
	{
		collection.GET("/:table_name", List)
		collection.POST("/:table_name", Create)
		collection.GET("/:table_name/:id", Read)
		collection.POST("/:table_name/:id", Update)
		collection.DELETE("/:table_name/:id", Delete)
	}

}
