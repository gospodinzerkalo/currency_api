package currency

import (
	"database/sql"
	"errors"
	"log"
	_ "github.com/lib/pq"
)

var currencyQueries = []string{
	`CREATE TABLE IF NOT EXISTS currency
	(
		title TEXT PRIMARY KEY,
		pub_date TEXT,
		description TEXT,
		quant FLOAT,
		index TEXT,
		change TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS history
	(
		id SERIAL PRIMARY KEY,
		title TEXT,
		pub_date TEXT,
		CHANGE TEXT
	);
	`,
}

type currencyStore struct {
	db *sql.DB
}

func NewCurrencyStore(cfg Config) (Repository, error) {
	dbConnection, err := getDbConn(getConnString(cfg))
	if err != nil {
		return nil, err
	}

	if err = dbConnection.Ping(); err != nil {
		return nil, err
	}

	for _, q := range currencyQueries {
		_, err = dbConnection.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}

	return currencyStore{dbConnection}, nil
}

func (store currencyStore) Create(currency XmlCurrency) (*XmlCurrency, error) {
	res, err := store.db.Exec("INSERT INTO currency (title, pub_date, description, quant, index, change)" +
		" VALUES ($1, $2, $3, $4, $5, $6)",
			currency.Title,
			currency.PubDate,
			currency.Description,
			currency.Quant,
			currency.Index,
			currency.Change,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if n <= 0 {
		return nil, errors.New("not able to create currency")
	}

	return &currency, nil
}

func (store currencyStore) Get(params *Params) ([]JsonCurrency, error) {
	var res []JsonCurrency
	var rows *sql.Rows
	var err error


	if params != nil && params.Base != "" && params.Quoted != "" {
		rows, err = store.db.Query("SELECT * FROM currency WHERE title=$1 OR title=$2", params.Base, params.Quoted)
	}else {
		rows, err = store.db.Query("SELECT * FROM currency")
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var curr JsonCurrency
		err := rows.Scan(&curr.Title, &curr.PubDate, &curr.Description, &curr.Quant, &curr.Index, &curr.Change)
		if err != nil {
			return nil, err
		}

		res = append(res, curr)
	}

	return res, nil
}

func (store currencyStore) Update(currency XmlCurrency) (*XmlCurrency, error) {
	res, err := store.db.Exec("UPDATE currency SET title=$1, pub_date=$2, description=$3, quant=$4, index=$5, change=$6 WHERE title=$7",
		currency.Title,
		currency.PubDate,
		currency.Description,
		currency.Quant,
		currency.Index,
		currency.Change,
		currency.Title,
		)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err

	}

	if n <= 0 {
		return nil, errors.New("not able to update currency with title="+currency.Title)
	}

	return &currency, nil
}

func (store currencyStore) CreateHistory(history History) (*History, error) {
	res, err := store.db.Exec("INSERT INTO history (title, pub_date, change) VALUES ($1, $2, $3)",
		history.Title,
		history.PubDate,
		history.Change,
		)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if n <= 0 {
		return nil, errors.New("not able to create history")
	}

	return &history, nil
}

func (store currencyStore) GetHistoryList() ([]History, error) {
	rows, err := store.db.Query("SELECT * FROM history")
	if err != nil {
		return nil, err
	}

	var res []History

	for rows.Next() {
		history := History{}
		if err := rows.Scan(&history.ID, &history.Title, &history.PubDate, &history.Change); err != nil {
			return nil, err
		}

		res = append(res, history)
	}

	return res, nil
}