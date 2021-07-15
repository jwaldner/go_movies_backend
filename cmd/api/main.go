package main

import "flag"

type config struct {
port int
env string	
}

func main () {

	var cfg config

	flag.IntVar(&cfg.port,"port",4000,"Server port to listen on")
	flag.StringVar(&cfg.env,"env","development","Application environment (development/production)")
	


}