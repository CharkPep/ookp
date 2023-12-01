package main

import (
	"fmt"
	"go.uber.org/fx"
	"ookp/src/controller/router"
	"ookp/src/model"
	"ookp/src/providers/datasource"
	firebase2 "ookp/src/providers/firebase"
)

func main() {

	app := fx.New(
		fx.Provide(model.Init),
		datasource.Module,
		firebase2.Module,
		router.ControllerRoutes,
		fx.Invoke(func(app *model.App) {
			fmt.Println("App started")
		}),
	)

	app.Run()

}
