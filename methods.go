package gaarx

type (
	methods struct {
		m map[string]func(app *App) error
	}
	Method struct {
		Name string
		Func func(app *App) error
	}
)
