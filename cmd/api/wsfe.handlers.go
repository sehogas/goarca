package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
	"github.com/sehogas/goarca/internal/util/validator"
)

// FEDummyHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200	{object}	dto.DummyResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/fe/Dummy [get]
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
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			ptoVta		query		string	true	"Punto de venta"
//	@Param			cbteTipo	query		string	true	"Tipo de comprobante"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/CompUltimoAutorizado [get]
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

// FECAESolicitarHandler godoc
//
//	@Summary		Solicitar CAE
//	@Description	Solicitar CAE
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Param			request		body		dto.FECAESolicitarRequest	true	"FECAESolicitarRequest"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/CAESolicitar [post]
func FECAESolicitarHandler(w http.ResponseWriter, r *http.Request) {
	var post dto.FECAESolicitarRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	resultado, err := Wsfe.FECAESolicitar(post.Cab, post.Det)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposCbteHandler godoc
//
//	@Summary		Listar Tipos de Comprobante
//	@Description	Este método permite consultar los tipos de comprobantes habilitados en este WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.CbteTipo
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposCbte [get]
func FEParamGetTiposCbteHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposCbte()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposConceptoHandler godoc
//
//	@Summary		Listar Tipos de Concepto
//	@Description	Este método devuelve los tipos de conceptos posibles en este WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.ConceptoTipo
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposConcepto [get]
func FEParamGetTiposConceptoHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposConcepto()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposDocHandler godoc
//
//	@Summary		Listar Tipos de Documento
//	@Description	Este método retorna el universo de tipos de documentos disponibles en el presente WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.DocTipo
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposDoc [get]
func FEParamGetTiposDocHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposDoc()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposIvaHandler godoc
//
// @Summary		Listar Tipos de IVA
// @Description	Mediante este método se obtiene la totalidad de alícuotas de IVA posibles de uso en el presente WS, detallando código y descripción.
// @Tags			Factura Electrónica
// @Produce		json
// @Param			x-api-key	header		string				true	"API Key de acceso"
// @Success		200			array		wsfe.IvaTipo
// @Failure		401			{object}	dto.ErrorResponse
// @Failure		500			{object}	dto.ErrorResponse
// @Router			/fe/GetTiposIva [get]
func FEParamGetTiposIvaHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposIva()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposMonedasHandler godoc
//
//	@Summary		Listar Tipos de Monedas
//	@Description	Este método retorna el universo de Monedas disponibles en el presente WS, indicando id y descripción de cada una.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.Moneda
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposMonedas [get]
func FEParamGetTiposMonedasHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposMonedas()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposOpcionalHandler godoc
//
//	@Summary		Listar Tipos Opcional
//	@Description	Este método permite consultar los códigos y descripciones de los tipos de datos Opcionales que se encuentran habilitados para ser usados en el WS
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.OpcionalTipo
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposOpcional [get]
func FEParamGetTiposOpcionalHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposOpcional()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetTiposTributosHandler godoc
//
//	@Summary		Listar Tipos Tributos
//	@Description	Devuelve los posibles códigos de tributos que puede contener un comprobante y su descripción.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.TributoTipo
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetTiposTributos [get]
func FEParamGetTiposTributosHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposTributos()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetPtosVentaHandler godoc
//
//	@Summary		Listar Puntos de Venta
//	@Description	Este método permite consultar los puntos de venta para ambos tipos de Código de Autorización (CAE y CAEA) gestionados previamente por la CUIT emisora.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Success		200			array		wsfe.PtoVenta
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/GetPtosVenta [get]
func FEParamGetPtosVentaHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetPtosVenta()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}
