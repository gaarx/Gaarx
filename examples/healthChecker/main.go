package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gaarx/gaarx"
	database "github.com/gaarx/gaarxDatabase"
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
	var conf conf.Config
	err := application.LoadConfig(*configSource, *configFile, &conf)
	if err != nil {
		panic(err)
	}
	err = application.InitializeLogger("file", conf.Log, &logrus.TextFormatter{
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
		database.WithDatabase(conf.DB, &entities.Resource{}, &entities.History{}),
	)
	go func() {
		sig := <-stop
		time.Sleep(2 * time.Second)
		finish()
		fmt.Printf("caught sig: %+v\n", sig)
		done <- true
	}()
	application.Start()
	<-ctx.Done()
	os.Exit(0)
}
