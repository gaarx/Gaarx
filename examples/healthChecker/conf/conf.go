package conf

type Config struct {
	DB   string
	Log  string
	Http struct {
		Addr string
	}
	Checker struct {
		Periodicity int
	}
}
