package main
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}