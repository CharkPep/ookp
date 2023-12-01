package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"ookp/src/model"
)

type MapRouter struct{}

func mapRoutes(app *model.App) MapRouter {
	app.Router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "OOKP assignment",
		})
	})

	return MapRouter{}
}

var ControllerRoutes = fx.Module("controller-routes",
	fx.Provide(mapRoutes),
	fx.Invoke(func(router MapRouter) {
		fmt.Printf("Routes mapped\n")
	}),
)
