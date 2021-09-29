package main

/*
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
        pql.conf.Postgres.User,
        pql.conf.Postgres.Password,
        pql.conf.Postgres.Host,
        pql.conf.Postgres.Port,
        pql.conf.Postgres.DB)
		*/
	
type config struct {
	port int
	env  string

	db struct {
		user string
		password string
		host string
		port string
		defaultDb string
		dsn string
	}
	jwt struct {
		secret string
	}

	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbData     string
}
