package services

import "github.com/gaarx/gaarx"

type CheckService struct {
	app *gaarx.App
}

func GetCheckService() gaarx.Service {
	return &CheckService{}
}

func (cs *CheckService) Start(app *gaarx.App) error {
	cs.app = app
	return nil
}

func (cs *CheckService) Stop() {

}

func (cs *CheckService) GetName() string {
	return "CheckService"
}
