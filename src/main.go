package main

import (
	"fmt"
	"go.uber.org/fx"
	"ookp/src/controller/router"
	"ookp/src/model"
)

func main() {

	app := fx.New(
		fx.Provide(model.Init),
		router.ControllerRoutes,
		fx.Invoke(func(app *model.App) {
			fmt.Println("App started")
		}),
	)

	app.Run()

}
