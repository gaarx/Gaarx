package gaarx

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	App struct {
		config
		database interface{}
		services map[string]Service
		ready    bool
		ctx      context.Context
		Finish   func()
		meta     *Meta
		storage  *storage
		methods  *methods
		events   map[string]*Event
		debug    bool
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

// Work Alias for Start
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
			log.Info().Str("service Name", n).Msg("Starting")
			err := s.Start(app)
			if err != nil {
				log.Warn().Err(err)
			}
		}(name, service)
	}
	for {
		select {
		case <-app.ctx.Done():
			log.Debug().Msg("Go to stop all")
			for _, service := range app.services {
				service.Stop()
			}
			log.Info().Msg("application was stopped")
			return
		default:
			runtime.Gosched()
			time.Sleep(time.Second)
		}
	}
}

// Stop all services and application
func (app *App) Stop() {
	log.Info().Msg("Application will be stopped")
	if app.Finish != nil {
		app.Finish()
	}
}

// GetService Return service and flag provide that service isset or not
func (app *App) GetService(name string) (Service, bool) {
	if service, ok := app.services[name]; ok {
		return service, true
	}
	return nil, false
}

// GetMeta Return Meta information
func (app *App) GetMeta() *Meta {
	return app.meta
}

// Storage Return all of application storage
func (app *App) Storage() *storage {
	return app.storage
}

// CallMethod call concrete method (should be defined in Initialize.WithMethods)
func (app *App) CallMethod(name string) error {
	if app.methods == nil {
		return errors.New("app run without methods option")
	}
	if function, ok := app.methods.m[name]; ok {
		return function(app)
	}
	return errors.New("app has no method " + name)
}

// Event Return concrete event
func (app *App) Event(name string) *Event {
	if event, ok := app.events[name]; ok {
		return event
	}
	ev := newEvent(name, app.ctx)
	ev.debug = app.debug
	app.events[name] = ev
	go ev.iterate()
	if app.debug {
		log.Debug().Msgf("Event %s was created", name)
	}
	return ev
}

// LoadConfig Load configuration from file
func (app *App) LoadConfig(configFile, configType string, configStruct interface{}) error {
	app.initConfig(configStruct)
	return app.loadConfig(configFile, configType)
}
