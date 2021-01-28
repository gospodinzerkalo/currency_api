package currency

import (
	"database/sql"
	"strconv"
)

type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	Params           string
	ConnectionString string
}


func getConnString(cfg Config) string {
	var connStr string
	if cfg.ConnectionString == "" {
		connStr = "postgres://"
		if cfg.Host == "" {
			cfg.Host = "localhost"
		}
		if cfg.Port == 0 {
			cfg.Port = 5432
		}
		connStr +=
			cfg.User + ":" +
				cfg.Password + "@" +
				cfg.Host + ":" +
				strconv.Itoa(cfg.Port) + "/" +
				cfg.Database + "?" +
				cfg.Params
	} else {
		connStr = cfg.ConnectionString
	}
	return connStr
}

func getDbConn(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr + " sslmode=disable")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
