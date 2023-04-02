package app

import (
	"log"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

type App struct {
	Config Environment
}

type Environment struct {
	Address          string `env:"RUN_ADDRESS" envDefault:":8080"`
	ConnectionString string `env:"DATABASE_URI"`
	SecretKey        string `env:"SECRET_KEY"`
}

func NewApp() App {
	app := App{Config: Environment{}}
	err := env.Parse(&app.Config)
	if err != nil {
		panic("can not create Application")
	}
	a := flag.StringP("Address", "a", "", "host for server")
	d := flag.StringP("DATABASE_URI", "d", "", "host for db")
	s := flag.StringP("SecretKey", "s", "", "secret Key")
	flag.Parse()
	if *a != "" {
		app.Config.Address = *a
	}
	if *d != "" {
		app.Config.ConnectionString = *d
	}

	if *s != "" {
		app.Config.SecretKey = *s
	}

	if app.Config.SecretKey == "" {
		app.Config.SecretKey = "secret"
	}

	if app.Config.ConnectionString == "" {
		panic("connection string is requered")
	}
	log.Printf("connection string is %s\n", app.Config.ConnectionString)
	return app
}
