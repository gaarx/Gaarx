package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gaarx/gaarx"
	database "github.com/gaarx/gaarxDatabase"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"healthchecker/conf"
	"healthchecker/entities"
	"healthchecker/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile = flag.String("config", "config.toml", "Project config")
var configSource = flag.String("source", "toml", "Project config source")

func main() {
	flag.Parse()
	var stop = make(chan os.Signal)
	var done = make(chan bool, 1)
	ctx, finish := context.WithCancel(context.Background())
	var application = &gaarx.App{}
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var config conf.Config
	err := application.LoadConfig(*configFile, *configSource, &config)
	if err != nil {
		panic(err)
	}
	err = application.InitializeLogger("file", config.Log, &logrus.TextFormatter{
		FullTimestamp: true,
	})
	if err != nil {
		panic(err)
	}
	application.Initialize(
		gaarx.WithContext(ctx),
		gaarx.WithServices(
			services.GetCheckService(),
			services.GetHttpService(ctx),
		),
		database.WithDatabase(
			config.DB,
			&entities.Resource{},
			&entities.History{},
		),
		gaarx.WithStorage(
			entities.ScopeResources,
			entities.ScopeHistories,
		),
		gaarx.WithMethods(
			gaarx.Method{
				Name: "GetResources",
				Func: func(app *gaarx.App) error {
					_ = app.Storage().ClearScope(entities.ScopeResources)
					var resources []*entities.Resource
					app.GetDB().(*gorm.DB).Find(&resources)
					for _, rs := range resources {
						err = app.Storage().Set(entities.ScopeResources, rs.Url, rs)
						if err != nil {
							app.GetLog().Error(err)
						}
					}
					return nil
				},
			},
		),
	)
	go func() {
		sig := <-stop
		time.Sleep(2 * time.Second)
		finish()
		fmt.Printf("caught sig: %+v\n", sig)
		done <- true
	}()
	_ = application.CallMethod("GetResources")
	application.Start()
	<-ctx.Done()
	os.Exit(0)
}
