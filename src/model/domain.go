package model

import (
	"context"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"path/filepath"
	"strings"
)

type App struct {
	Router *gin.Engine
	Server *http.Server
}

func (Self *App) mapViews() {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob("./src/view/layouts/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob("./src/view/includes/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	fmt.Printf("layouts: %v\n", layouts)
	fmt.Printf("includes: %v\n", includes)
	for _, include := range includes {
		includeName := filepath.Base(filepath.Dir(include))
		fmt.Printf("includeName: %v\n", includeName)
		includeLayouts := []string{}
		for _, layout := range layouts {
			if strings.Contains(layout, includeName) {
				includeLayouts = append(includeLayouts, layout)
			}
		}
		fmt.Printf("includeLayouts: %v\n", includeLayouts)
		fmt.Printf("include: %v\n", include)
		r.AddFromFiles(filepath.Base(include), append(includeLayouts, include)...)

	}
	Self.Router.HTMLRender = r
}

func Init(lc fx.Lifecycle) *App {
	router := gin.Default()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	app := &App{
		Router: router,
		Server: server,
	}
	app.mapViews()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			router.Static("/public", "./public")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return app
}
