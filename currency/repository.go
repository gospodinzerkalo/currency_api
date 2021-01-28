package currency

type Repository interface {
	Create(currency XmlCurrency) (*XmlCurrency, error)
	Get(params *Params) ([]JsonCurrency, error)
	Update(currency XmlCurrency) (*XmlCurrency, error)
	CreateHistory(history History) (*History, error)
	GetHistoryList() ([]History, error)
}
