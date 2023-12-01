package datasource

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"ookp/src/model"
	"os"
	"time"
)

type Datasource struct {
	DB *gorm.DB
}

func NewDatasource(lc fx.Lifecycle) *Datasource {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432"
	var db *gorm.DB
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: dbLogger,
			}); err != nil {
				panic(err)
			}

			if err = db.AutoMigrate(&model.User{}); err != nil {
				panic(err)
			}
			return nil
		},
	})

	return &Datasource{
		DB: db,
	}
}

var Module = fx.Module("Datasource",
	fx.Provide(NewDatasource),
	fx.Invoke(func(datasource *Datasource) {
		fmt.Printf("Datasource: %v\n", datasource)
	}))
