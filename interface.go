package gaarx

type (
	ConfigAble interface {
		GetConnString() string
		GetLogWay() string
		GetLogDestination() string
		GetLogApplicationName() string
	}
)
