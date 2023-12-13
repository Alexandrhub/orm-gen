package utils

// DB структура базы данных
type DB struct {
	Net      string
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	MaxConn  int
	Port     string
	Timeout  int
}
