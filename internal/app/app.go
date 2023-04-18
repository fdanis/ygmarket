package app

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env"
	"github.com/fdanis/yg-loyalsys/internal/db/driver"
	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	flag "github.com/spf13/pflag"
)

type App struct {
	Config             Environment
	db                 *sql.DB
	OrderRepository    repositories.OrderRepository
	WithdrawRepository repositories.WithdrawRepository
	UserRepository     repositories.UserRepository
}

type Environment struct {
	Address              string `env:"RUN_ADDRESS" envDefault:":8080"`
	ConnectionString     string `env:"DATABASE_URI"`
	SecretKey            string `env:"SECRET_KEY"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
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
	r := flag.StringP("ACCRUAL_SYSTEM_ADDRESS", "r", "", "ACCRUAL_SYSTEM_ADDRESS")
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

	if *r != "" {
		app.Config.AccrualSystemAddress = *r
	}

	if app.Config.SecretKey == "" {
		app.Config.SecretKey = "secret"
	}

	if app.Config.ConnectionString == "" {
		panic("connection string is requered")
	}

	if app.Config.AccrualSystemAddress == "" {
		panic("AccrualSystemAddress is requered")
	}

	log.Printf("connection string is %s\n", app.Config.ConnectionString)

	db, err := driver.ConnectSQL(app.Config.ConnectionString)
	if err != nil {
		panic(err)
	}
	app.db = db
	app.OrderRepository = repositories.NewOrderRepository(db)
	app.UserRepository = repositories.NewUserRepository(db)
	app.WithdrawRepository = repositories.NewWithdrawRepository(db)

	return app
}

func (a *App) Close() {
	a.db.Close()
}
