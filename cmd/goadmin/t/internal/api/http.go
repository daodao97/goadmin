package api

import (
	"admin/internal/conf"
	"github.com/gin-gonic/gin"
)

func Route(e *gin.Engine, conf *conf.Conf) {
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
