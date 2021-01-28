package currency

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Service interface {
	RefreshCurrencies() error
	ParseXml() (*[]XmlCurrency, error)
	GetCurrencies(params Params) (*[]ResponseCurrency, error)
	Convert(convert Convert)  (*ConvertResponse, error)
	GetHistoryList() ([]History, error)
}

const (
	XML_SOURCE = "https://nationalbank.kz/rss/rates_all.xml?switch=kazak"

)

type ServiceCurrency struct {
	Repository 		Repository
}

func NewService(repository Repository) Service {
	return &ServiceCurrency{Repository: repository}
}

func (service *ServiceCurrency) RefreshCurrencies() error {
	res, err := service.ParseXml()
	if err != nil {
		return err
	}
	currencies, err := service.Repository.Get(nil)
	if err != nil {
		return err
	}
	n := len(currencies)
	for _, curr := range *res {
		var er error
		switch n {
		case 0:
			_, er = service.Repository.Create(curr)
		default:
			_, er = service.Repository.Update(curr)

		}
		if er != nil {
			return nil
		}
		if _, err := service.Repository.CreateHistory(History{
			Title:   curr.Title,
			PubDate: curr.PubDate,
			Change:  curr.Change,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (service *ServiceCurrency) ParseXml() (*[]XmlCurrency, error) {
	resp, err := http.Get(XML_SOURCE)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var curr XmlCurrencyResp
	if err := xml.Unmarshal(data, &curr); err != nil {
		return nil, err
	}

	return &curr.Currencies, nil
}

func (service *ServiceCurrency) GetCurrencies(params Params) (*[]ResponseCurrency, error) {
	res, err := service.Repository.Get(&params)
	if err != nil {
		return nil, err
	}

	var resp []ResponseCurrency

	if params.Base != "" && params.Quoted != "" && len(res) == 2 {
		pairCurr := ResponseCurrency{}
		pairCurr.CurrencyPair = params.Base + "/" + params.Quoted
		switch res[0].Title {
		case params.Base:
			pairCurr.Value = fmt.Sprintf("%.2f", res[0].Description / res[1].Description)
		case params.Quoted:
			pairCurr.Value = fmt.Sprintf("%.2f", res[1].Description / res[0].Description)
		}
		resp = append(resp, pairCurr)
	} else {
		for _, curr := range res {
			pairCurr := ResponseCurrency{
				CurrencyPair: curr.Title+ "/KZT",
				Value:        fmt.Sprintf("%v", curr.Description),
			}

			resp = append(resp, pairCurr)
		}
	}

	return &resp, nil
}

func (service *ServiceCurrency) Convert(convert Convert) (*ConvertResponse, error) {
	params := Params{
		Base:   convert.ConvertFrom,
		Quoted: convert.ConvertTo,
	}

	getCurr, err := service.Repository.Get(&params)
	if err != nil {
		return nil, err
	}

	res := ConvertResponse{}
	res.ConvertFrom = convert.ConvertFrom
	res.ConvertTo = convert.ConvertTo
	res.ConvertFromValue = convert.Value

	if len(getCurr) == 2 {
		switch getCurr[0].Title {
		case params.Base:
			res.ConvertToValue = fmt.Sprintf("%.4f", (getCurr[0].Description / float64(getCurr[0].Quant)) / (getCurr[1].Description * convert.Value))
		case params.Quoted:
			res.ConvertToValue = fmt.Sprintf("%.4f", (getCurr[1].Description / float64(getCurr[1].Quant)) / (getCurr[0].Description * convert.Value))
		}
	}else {
		switch params.Base {
		case "KZT":
			res.ConvertToValue = fmt.Sprintf("%.4f", convert.Value / getCurr[0].Description)
		default:
			res.ConvertToValue = fmt.Sprint(getCurr[0].Description)

		}
	}


	return &res, nil
}

func (service *ServiceCurrency) GetHistoryList() ([]History, error) {
	return service.Repository.GetHistoryList()
}