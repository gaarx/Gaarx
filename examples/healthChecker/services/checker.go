package services

import (
	"context"
	"fmt"
	"github.com/gaarx/gaarx"
	"github.com/jinzhu/gorm"
	"healthchecker/conf"
	"healthchecker/entities"
	"net/http"
	"time"
)

type CheckService struct {
	resources []*entities.Resource

	app              *gaarx.App
	ctx              context.Context
	cwork            context.Context
	cancel           context.CancelFunc
	wcancel          context.CancelFunc
	allDone          chan interface{}
	checkPeriodicity time.Duration
}

func GetCheckService() gaarx.Service {
	return &CheckService{}
}

func (cs *CheckService) Start(app *gaarx.App) error {
	cs.app = app
	periodicity := app.Config().(*conf.Config).Checker.Periodicity
	if periodicity > 0 {
		cs.checkPeriodicity = time.Second * time.Duration(periodicity)
	} else {
		cs.checkPeriodicity = time.Minute // Default time
	}
	cs.ctx, cs.cancel = context.WithCancel(context.Background())
	cs.cwork, cs.wcancel = context.WithCancel(context.Background())
	cs.allDone = make(chan interface{})
	cs.loadResources()
	go cs.work()
	go func() {
		for {
			select {
			case _ = <-cs.app.Event(entities.EventReloadResources).Listen():
				_ = app.CallMethod("GetResources")
				cs.loadResources()
				go cs.work()
			case <-cs.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (cs *CheckService) work() {
	t := time.NewTicker(cs.checkPeriodicity)
	for {
		select {
		case <-t.C:
			for _, resource := range cs.resources {
				start := time.Now()
				resp, err := http.Get(fmt.Sprintf("http://%s", resource.Url))
				end := time.Now()
				h := entities.History{
					Url:         resource.Url,
					CheckTime:   time.Now(),
					RequestTime: int(end.Sub(start).Milliseconds()),
				}
				if err != nil {
					cs.app.GetLog().Error(err)
					h.Status = fmt.Sprintf("Request failed: %v", err.Error())
					cs.app.GetDB().(*gorm.DB).Save(&h)
					continue
				}
				h.StatusCode = resp.StatusCode
				h.Status = "Success"
				cs.app.GetDB().(*gorm.DB).Save(&h)
			}
		case <-cs.ctx.Done():
			cs.allDone <- true
			return
		case <-cs.cwork.Done():
			return
		}
	}
}

func (cs *CheckService) Stop() {
	cs.cancel()
	<-cs.allDone
}

func (cs *CheckService) GetName() string {
	return "CheckService"
}

func (cs *CheckService) loadResources() {
	rs, err := cs.app.Storage().GetAll(entities.ScopeResources)
	if err != nil {
		panic(err)
	}
	cs.resources = make([]*entities.Resource, 0)
	for _, r := range rs {
		cs.resources = append(cs.resources, r.(*entities.Resource))
	}
}
