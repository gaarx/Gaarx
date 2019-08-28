package gaarx

import (
	"context"
	"errors"
	"fmt"
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

type (
	App struct {
		config
		database interface{}
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

// Initialize application with option. Logger and Database should be initialized yet
func (app *App) Initialize(options ...Option) *App {
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
	app.ready = true
	return app
}

// Alias for Start
func (app *App) Work() {
	app.Start()
}

// Start application, all services and wait for calling Stop() method
func (app *App) Start() {
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

// Stop all services and application
func (app *App) Stop() {
	app.GetLog().Info("Application will be stopped")
	if app.Finish != nil {
		app.Finish()
	}
}

// Return logrus instance
func (app *App) GetLog() *logrus.Logger {
	return app.log
}

// Return service and flag provide that service isset or not
func (app *App) GetService(name string) (Service, bool) {
	if service, ok := app.services[name]; ok {
		return service, true
	}
	return nil, false
}

// Return Meta information
func (app *App) GetMeta() *Meta {
	return app.meta
}

// Return all of application storage
func (app *App) Storage() *storage {
	return app.storage
}

// Call concrete method (should be defined in Initialize.WithMethods)
func (app *App) CallMethod(name string) error {
	if app.methods == nil {
		return errors.New("app run without methods option")
	}
	if function, ok := app.methods.m[name]; ok {
		return function(app)
	}
	return errors.New("app has no method " + name)
}

// Return concrete event
func (app *App) Event(name string) *Event {
	if event, ok := app.events[name]; ok {
		return event
	}
	ev := newEvent(name, app.ctx)
	app.events[name] = ev
	go ev.iterate()
	return ev
}

// Load configuration from file
func (app *App) LoadConfig(way string, configFile string, configStruct interface{}) error {
	app.initConfig(configStruct)
	return app.loadConfig(ConfigSource(way), configFile)
}

// Initialize logger with given settings
func (app *App) InitializeLogger(way LogWay, path string, formatter logrus.Formatter) error {
	app.log = logrus.New()
	switch way {
	case FileLog:
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		app.log.Out = f
		break
	case GrayLog:
		app.log.AddHook(graylog.NewGraylogHook(path, map[string]interface{}{}))
		break
	default:
		app.log.Out = os.Stdout
	}
	app.log.SetFormatter(formatter)
	return nil
}
