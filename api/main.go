package main

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shellhub-io/shellhub/api/store/mongo"
	"github.com/shellhub-io/shellhub/api/v2/app"
	"github.com/shellhub-io/shellhub/api/v2/handlers"
	mongostore "github.com/shellhub-io/shellhub/api/v2/store/mongo"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	MongoHost string `envconfig:"mongo_host" default:"mongo"`
	MongoPort int    `envconfig:"mongo_port" default:"27017"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	var cfg config
	if err := envconfig.Process("api", &cfg); err != nil {
		panic(err.Error())
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort))
	// Connect to MongoDB
	client, err := mgo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	if err := mongo.ApplyMigrations(client.Database("main")); err != nil {
		panic(err)
	}

	/*e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, latecy=${latency_human}\n",
	}))*/

	mongoStore := mongostore.NewStore(client.Database("main"))

	e.GET("/merda", func(c echo.Context) error {
		return nil
	})

	handlers.NewHandler(app.NewApp(mongoStore), e)

	e.Logger.Fatal(e.Start(":8080"))
}
