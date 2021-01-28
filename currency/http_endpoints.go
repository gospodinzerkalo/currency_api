package currency

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpEndpointFactory interface {
	MakeGetCurrency() func(w http.ResponseWriter, r *http.Request)
	MakeConvert() func(w http.ResponseWriter, r *http.Request)
	MakeGetHistoryList() func(w http.ResponseWriter, r *http.Request)
}

type httpEndpointFactory struct {
	service Service
}

func NewHttpEndpointFactory(service Service) HttpEndpointFactory {
	return httpEndpointFactory{service: service}
}

func (httpFac httpEndpointFactory) MakeGetCurrency() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		base := r.URL.Query().Get("base")
		quoted := r.URL.Query().Get("quoted")
		params := Params{
			Base:   base,
			Quoted: quoted,
		}
		resp, err := httpFac.service.GetCurrencies(params)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}

		response, err := json.Marshal(resp)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		writeResponse(w, http.StatusOK, response)
	}
}

func (httpFac httpEndpointFactory) MakeConvert() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		converter := Convert{}

		if err := json.Unmarshal(body, &converter); err != nil {
			writeResponse(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		resp, err := httpFac.service.Convert(converter)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		response, err := json.Marshal(resp)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		writeResponse(w, http.StatusOK, response)
	}
}

func (httpFac httpEndpointFactory) MakeGetHistoryList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res, err := httpFac.service.GetHistoryList()
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		resp, err := json.Marshal(res)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		writeResponse(w, http.StatusOK, resp)
	}
}

func writeResponse(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
}