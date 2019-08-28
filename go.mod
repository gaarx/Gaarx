module github.com/gaarx/gaarx

go 1.12

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20181115231424-8e868ca12c0f

replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1

replace layeh.com/radius => github.com/layeh/radius v0.0.0-20190118135028-0f678f039617

require (
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/gemnasium/logrus-graylog-hook v2.0.7+incompatible // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.4
	github.com/micro/go-config v1.1.0
	github.com/sirupsen/logrus v1.4.1
)
