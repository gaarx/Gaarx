package gaarx

import (
	"context"
	"sync"
)

type (
	Option func(app *App) error
)

func WithServices(services ...Service) Option {
	return func(a *App) error {
		a.services = make(map[string]Service)
		for _, s := range services {
			a.services[s.GetName()] = s
		}
		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(a *App) error {
		a.ctx = ctx
		return nil
	}
}

func WithMetaInformation(version, build, buildTime, branch, commit string) Option {
	return func(a *App) error {
		a.meta = &Meta{
			Version:   version,
			Build:     build,
			BuildTime: buildTime,
			Branch:    branch,
			Commit:    commit,
		}
		return nil
	}
}

func WithStorage(scopes ...string) Option {
	return func(a *App) error {
		a.storage = &storage{
			innerMap: make(map[string]*sync.Map),
		}
		for _, scope := range scopes {
			a.storage.innerMap[scope] = &sync.Map{}
		}
		return nil
	}
}

func WithMethods(functions ...Method) Option {
	return func(a *App) error {
		a.methods = &methods{
			m: make(map[string]func(*App) error),
		}
		for _, function := range functions {
			a.methods.m[function.Name] = function.Func
		}
		return nil
	}
}
