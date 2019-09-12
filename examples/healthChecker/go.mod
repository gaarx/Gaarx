module healthchecker

go 1.12

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20181115231424-8e868ca12c0f

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43

require (
	github.com/gaarx/gaarx v0.1.0
	github.com/gaarx/gaarxDatabase v0.0.0-20190828060522-63aad25d99e4
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/sirupsen/logrus v1.4.1
	github.com/ugorji/go v1.1.7 // indirect
)
