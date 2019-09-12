package services

import (
	"context"
	"github.com/gaarx/gaarx"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"healthchecker/conf"
	"healthchecker/entities"
	"net/http"
)

type HttpService struct {
	app *gaarx.App
	ctx context.Context
	srv *http.Server
}

func GetHttpService(ctx context.Context) gaarx.Service {
	return &HttpService{
		ctx: ctx,
	}
}

func (hs *HttpService) Start(app *gaarx.App) error {
	hs.app = app
	r := gin.Default()
	conf := app.Config().(*conf.Config).Http
	srv := &http.Server{
		Addr:    conf.Addr,
		Handler: r,
	}
	private := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))
	private.GET("/", index)
	private.POST("/resource/add", hs.addResource())
	private.GET("/history", hs.getResourceHistory())
	hs.srv = srv
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (hs *HttpService) Stop() {
	if hs.srv != nil {
		_ = hs.srv.Shutdown(context.Background())
	}
}

func (hs *HttpService) GetName() string {
	return "HttpService"
}

func index(c *gin.Context) {
	c.HTML(200, "http", gin.H{
		"status": "success",
	})
}

func (hs *HttpService) addResource() func(c *gin.Context) {
	return func(c *gin.Context) {
		url := c.PostForm("url")
		db, ok := hs.app.GetDB().(*gorm.DB)
		if !ok {
			c.HTML(500, "error", gin.H{
				"status": "internal error",
			})
		}
		res := entities.Resource{Url: url}
		db.Create(&res)
		c.HTML(204, "success", gin.H{
			"status": "created",
			"id":     res.ID,
		})
	}
}

func (hs *HttpService) getResourceHistory() func(c *gin.Context) {
	return func(c *gin.Context) {
		url := c.Query("url")
		db, ok := hs.app.GetDB().(*gorm.DB)
		if !ok {
			c.HTML(500, "error", gin.H{
				"status": "internal error",
			})
		}
		var res []entities.History
		db.Where("url = ?", url).Find(&res)
		c.JSON(200, res)
	}
}
