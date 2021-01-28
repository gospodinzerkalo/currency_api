package currency

import "encoding/xml"

type JsonCurrency struct {
	Title 		string 		`json:"title"`
	PubDate 	string 		`json:"pub_date"`
	Description float64		`json:"description"`
	Quant		int			`json:"quant"`
	Index 		string		`json:"index"`
	Change 		string		`json:"change"`
}

type ResponseCurrency struct {
	CurrencyPair 	string 		`json:"currency_pair"`
	Value 			string		`json:"value"`
}

type XmlCurrency struct {
	Title 		string 		`xml:"title"`
	PubDate 	string		`xml:"pubDate"`
	Description	float64		`xml:"description"`
	Quant 		int 		`xml:"quant"`
	Index		string		`xml:"index"`
	Change 		string		`xml:"change"`
}

type XmlCurrencyResp struct {
	XMLName   	xml.Name 			`xml:"rss"`
	Currencies	[]XmlCurrency		`xml:"channel>item"`
}

type Params struct {
	Base 		string
	Quoted 		string
}

type Convert struct {
	ConvertFrom 	string 		`json:"convert_from"`
	ConvertTo		string		`json:"convert_to"`
	Value			float64		`json:"value"`
}

type ConvertResponse struct {
	ConvertFrom 			string 		`json:"convert_from"`
	ConvertTo				string		`json:"convert_to"`
	ConvertFromValue		float64		`json:"convert_from_value"`
	ConvertToValue			string		`json:"convert_to_value"`
}

type History struct {
	ID 			int64		`json:"id"`
	Title 		string		`json:"title"`
	PubDate 	string		`json:"pub_date"`
	Change 		string		`json:"change"`
}