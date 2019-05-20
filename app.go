package gaarx

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"os"
	"runtime"
	"time"
)

type (
	App struct {
		config
		db
		services map[string]Service
		ready    bool
		ctx      context.Context
		Finish   func()
		log      *logrus.Logger
		meta     *Meta
		storage  *storage
		methods  *methods
		events   map[string]*Event
	}
	Service interface {
		Start(app *App) error
		Stop()
		GetName() string
	}
)

func (app *App) Start(options ...Option) *App {
	for _, option := range options {
		e := option(app)
		if e != nil {
			fmt.Printf("Some option can't executed. See error: %s", e.Error())
			panic("Can't execute application option")
		}
	}
	if app.ctx == nil {
		ctx, finish := context.WithCancel(context.Background())
		app.ctx, app.Finish = ctx, finish
	}
	app.events = make(map[string]*Event)
	app.initializeLogger()
	app.initializeDatabase()
	app.ready = true
	return app
}

func (app *App) Work() {
	if !app.ready {
		panic("try to work not ready application")
	}
	for name, service := range app.services {
		go func(n string, s Service) {
			app.log.Info("Starting ", n)
			err := s.Start(app)
			if err != nil {
				app.log.Warning(err)
			}
		}(name, service)
	}
	for {
		select {
		case <-app.ctx.Done():
			app.log.Debug("Go to stop all")
			for _, service := range app.services {
				service.Stop()
			}
			app.log.Info("application was stopped")
			return
		default:
			runtime.Gosched()
			time.Sleep(time.Second)
		}
	}
}

func (app *App) Stop() {
	if app.Finish != nil {
		app.Finish()
	}
}

func (app *App) initializeLogger() {
	if app.configData.GetLogWay() != "" {
		f, err := os.OpenFile(app.configData.GetLogDestination(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		app.log.Out = f
	} else {
		panic("Logger config doesn't configured")
	}
}

func (app *App) initializeDatabase() {
	connString := app.Config().GetConnString()
	if connString == "" {
		return
		}
	db, err := gorm.Open(
		"mysql",
		app.Config().GetConnString(),
	)
	if err != nil {
		app.log.Fatal(err)
		panic(err)
	}
	app.database = db
	app.database.SetLogger(app.log)
	app.database.Set("gorm:table_options", "CHARSET=utf8")
	app.MigrateEntities(app.migrateEntities...)
}

func (app *App) GetLog() *logrus.Logger {
	return app.log
}

func (app *App) GetService(name string) (Service, bool) {
	if service, ok := app.services[name]; ok {
		return service, true
	}
	return nil, false
}

func (app *App) GetMeta() *Meta {
	return app.meta
}

func (app *App) Storage() *storage {
	return app.storage
}

func (app *App) CallMethod(name string) error {
	if app.methods == nil {
		return errors.New("app run without methods option")
	}
	if function, ok := app.methods.m[name]; ok {
		return function(app)
	}
	return errors.New("app has no method " + name)
}

func (app *App) Event(name string) *Event {
	if event, ok := app.events[name]; ok {
		return event
	}
	ev := newEvent(name, app.ctx)
	app.events[name] = ev
	go ev.iterate()
	return ev
}
