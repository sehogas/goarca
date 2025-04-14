package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

// DummyGesTabRefHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Success		200	{object}	dto.DummyResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/fe/dummy [get]
func FEDummyHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wsfe.FEDummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, &dto.DummyResponse{
		AppServer:  appServer,
		AuthServer: authServer,
		DbServer:   dbServer,
	}, nil)
}

// FECompUltimoAutorizadoHandler godoc
//
//	@Summary		Consultar último comprobante autorizado
//	@Description	Consultar último comprobante autorizado
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			ptoVta		query		string	true	"Punto de venta"
//	@Param			cbteTipo	query		string	true	"Tipo de comprobante"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/comp-ultimo-autorizado [get]
func FECompUltimoAutorizadoHandler(w http.ResponseWriter, r *http.Request) {
	ptoVtaStr := r.URL.Query().Get("ptoVta")
	cbteTipoStr := r.URL.Query().Get("cbteTipo")

	ptoVta, err := strconv.Atoi(ptoVtaStr)
	if err != nil {
		err := errors.New("error leyendo parámetro ptoVta")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	cbteTipo, err := strconv.Atoi(cbteTipoStr)
	if err != nil {
		err := errors.New("error leyendo parámetro cbteTipo")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cbteNro, err := Wsfe.FEUltimoComprobanteEmitido(int32(ptoVta), int32(cbteTipo))
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, &dto.CbteNroResponse{
		CbteNro: cbteNro,
	}, nil)
}
