package gaarx

func (app *App) SetDatabase(db interface{}) {
	app.database = db
}

func (app *App) GetDB() interface{} {
	return app.database
}
