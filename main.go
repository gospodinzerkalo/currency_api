package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gospodinzerkalo/currency_api/currency"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

var (
	configPath 				= "./"
	postgreDb				= ""
	postgreUser 			= ""
	postgrePass 			= ""
	postgreHost 			= ""
	postgrePort 			= 5432

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Usage:       "path to .env config file",
			Destination: &configPath,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "Currency Api"
	app.Usage = "origin run"
	app.UsageText = "origin run"

	app.Flags = flags
	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func parseEnv() {
	if configPath != "" {
		godotenv.Overload(configPath)
	}

	postgreHost = os.Getenv("POSTGRES_HOST")
	postgreDb = os.Getenv("POSTGRES_DATABASE")
	postgrePass = os.Getenv("POSTGRES_PASSWORD")
	postgreUser = os.Getenv("POSTGRES_USER")
}

func run(c *cli.Context) error {
	parseEnv()
	postgresConfig := currency.Config{
		Host:             postgreHost,
		Port:             postgrePort,
		User:             postgreUser,
		Password:         postgrePass,
		Database:         postgreDb,
		Params:           "",
		ConnectionString: "",
	}

	currencyStore, err := currency.NewCurrencyStore(postgresConfig)
	if err != nil {
		return err
	}

	currencyService := currency.NewService(currencyStore)
	currencyHttpFac := currency.NewHttpEndpointFactory(currencyService)

	router := mux.NewRouter()

	router.Methods("GET").Path("/currency").HandlerFunc(currencyHttpFac.MakeGetCurrency())
	router.Methods("POST").Path("/convert").HandlerFunc(currencyHttpFac.MakeConvert())
	router.Methods("GET").Path("/history").HandlerFunc(currencyHttpFac.MakeGetHistoryList())

	cr := cron.New()
	_, err = cr.AddFunc("@every 10s", func() {
		if err := currencyService.RefreshCurrencies(); err != nil {
			log.Fatal(err)
		}
	})

	cr.Start()

	if err != nil {
		return err
	}

	log.Println("Listen on 8080")
	log.Println(http.ListenAndServe("0.0.0.0:8080", router))

	return nil
}