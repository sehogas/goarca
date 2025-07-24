package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

// FEDummyHandler godoc
//
//	@Summary		Estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEDummy [get]
func FEDummyHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEDummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FECompUltimoAutorizadoHandler godoc
//
//	@Summary		Último comprobante autorizado
//	@Description	Consultar último comprobante autorizado
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			ptoVta		query		string	true	"Punto de venta"
//	@Param			cbteTipo	query		string	true	"Tipo de comprobante"
//	@Success		200			{object}	wsfe.FERecuperaLastCbteResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECompUltimoAutorizado [get]
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

	resultado, err := Wsfe.FEUltimoComprobanteEmitido(int32(ptoVta), int32(cbteTipo))
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FECAESolicitarHandler godoc
//
//	@Summary		Solicitar CAE
//	@Description	Solicitar CAE
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		dto.FECAESolicitarRequest	true	"FECAESolicitarRequest"
//	@Success		200			{object}	wsfe.FECAEResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECAESolicitar [post]
func FECAESolicitarHandler(w http.ResponseWriter, r *http.Request) {
	var post dto.FECAESolicitarRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
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
//	@Summary		Tipos de Comprobante
//	@Description	Este método permite consultar los tipos de comprobantes habilitados en este WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.CbteTipoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposCbte [get]
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
//	@Summary		Tipos de Concepto
//	@Description	Este método devuelve los tipos de conceptos posibles en este WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.ConceptoTipoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposConcepto [get]
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
//	@Summary		Tipos de Documento
//	@Description	Este método retorna el universo de tipos de documentos disponibles en el presente WS.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.DocTipoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposDoc [get]
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
//	@Summary		Tipos de IVA
//	@Description	Mediante este método se obtiene la totalidad de alícuotas de IVA posibles de uso en el presente WS, detallando código y descripción.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.IvaTipoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposIva [get]
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
//	@Summary		Tipos de Monedas
//	@Description	Este método retorna el universo de Monedas disponibles en el presente WS, indicando id y descripción de cada una.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.MonedaResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposMonedas [get]
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
//	@Summary		Tipos Opcional
//	@Description	Este método permite consultar los códigos y descripciones de los tipos de datos Opcionales que se encuentran habilitados para ser usados en el WS
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.OpcionalTipoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposOpcional [get]
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
//	@Summary		Tipos Tributos
//	@Description	Devuelve los posibles códigos de tributos que puede contener un comprobante y su descripción.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.FETributoResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposTributos [get]
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
//	@Summary		Puntos de Venta
//	@Description	Este método permite consultar los puntos de venta para ambos tipos de Código de Autorización (CAE y CAEA) gestionados previamente por la CUIT emisora.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.FEPtoVentaResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetPtosVenta [get]
func FEParamGetPtosVentaHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetPtosVenta()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetCotizacionHandler godoc
//
//	@Summary		Cotización de moneda
//	@Description	Retorna la última cotización de la base de datos aduanera de la moneda ingresada. Este valor es orientativo.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			monId		query		string	true	"código de moneda"
//	@Param			fchCotiz	query		string	false	"fecha de cotización"
//	@Success		200			{object}	wsfe.FECotizacionResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetCotizacion [get]
func FEParamGetCotizacionHandler(w http.ResponseWriter, r *http.Request) {
	monIdStr := r.URL.Query().Get("monId")
	fchCotiz := r.URL.Query().Get("fchCotiz")

	if monIdStr == "" {
		err := errors.New("error leyendo parámetros monId")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cotizacion, err := Wsfe.FEParamGetCotizacion(monIdStr, fchCotiz)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, cotizacion, nil)
}

// FECompTotXRequestHandler godoc
//
//	@Summary		Cantidad máxima de registros por request
//	@Description	Retorna la cantidad máxima de registros que se podrá incluir en un request al método FECAESolicitar / FECAEARegInformativo.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.FERegXReqResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECompTotXRequest [get]
func FECompTotXRequestHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FECompTotXRequest()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FECAEARegInformativo godoc
//
//	@Summary		Informar comprobantes emitidos y asociados a una CAEA
//	@Description	Este método permite informar para cada CAEA otorgado, la totalidad de los comprobantes emitidos y asociados a cada CAEA
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		dto.FeCAEARegInfReqRequest	true	"FeCAEARegInfReqRequest"
//	@Success		200			{object}	wsfe.FECAEAResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECAEARegInformativo [post]
func FECAEARegInformativoHandler(w http.ResponseWriter, r *http.Request) {
	var post dto.FeCAEARegInfReqRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wsfe.FECAEARegInformativo(post.Cab, post.Det)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FECAEASinMovimientoConsultarHandler godoc
//
//	@Summary		Consulta de Puntos Venta sin movimientos
//	@Description	Esta operación permite consultar mediante un CAEA, cuáles fueron los puntos de venta que fueron notificados como sin movimiento. El cliente envía el requerimiento, el cual es atendido por el WS, superadas las validaciones de seguridad se informa el CAEA, puntos de venta identificados como sin movimientos y fecha de proceso. En caso de informar el punto de venta, se informan los datos vinculados a ese punto de venta en particular.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			CAEA		query		string	true	"CAEA"
//	@Param			PtoVta		query		string	true	"Punto de Venta"
//	@Success		200			{object}	wsfe.FECotizacionResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECAEASinMovimientoConsultar [get]
func FECAEASinMovimientoConsultarHandler(w http.ResponseWriter, r *http.Request) {
	caea := r.URL.Query().Get("CAEA")
	ptoVtaStr := r.URL.Query().Get("PtoVta")

	if caea == "" {
		err := errors.New("error leyendo parámetros CAEA")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	ptoVta, err := strconv.Atoi(ptoVtaStr)
	if err != nil {
		err := errors.New("error leyendo parámetro PtoVta")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cotizacion, err := Wsfe.FECAEASinMovimientoConsultar(caea, int32(ptoVta))
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, cotizacion, nil)
}

// FECompConsultarHandler godoc
//
//	@Summary		Consulta datos de un Comprobante
//	@Description	Esta operación permite consultar mediante tipo, numero de comprobante y punto de venta los datos de un comprobante ya emitido. Dentro de los datos del comprobante resultante se obtiene el tipo de emisión utilizado para generar el código de autorización.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			PtoVta		query		string	true	"Punto de Venta"
//	@Param			CbteTipo	query		string	true	"Tipo de Comprobante"
//	@Param			CbteNro		query		string	true	"Número de Comprobante"
//	@Success		200			{object}	wsfe.FECompConsultaResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FECompConsultar [get]
func FECompConsultarHandler(w http.ResponseWriter, r *http.Request) {
	ptoVtaStr := r.URL.Query().Get("PtoVta")
	cbteTipoStr := r.URL.Query().Get("CbteTipo")
	cbteNroStr := r.URL.Query().Get("CbteNro")

	ptoVta, err := strconv.Atoi(ptoVtaStr)
	if err != nil {
		err := errors.New("error leyendo parámetro PtoVta")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cbteTipo, err := strconv.Atoi(cbteTipoStr)
	if err != nil {
		err := errors.New("error leyendo parámetro CbteTipo")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cbteNro, err := strconv.ParseInt(cbteNroStr, 10, 64)
	if err != nil {
		err := errors.New("error leyendo parámetro CbteNro")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	cotizacion, err := Wsfe.FECompConsultar(int32(ptoVta), int32(cbteTipo), cbteNro)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, cotizacion, nil)
}

// FEParamGetTiposPaisesHandler godoc
//
//	@Summary		Tipos de Países
//	@Description	Esta operación permite consultar los códigos de países y descripción de los mismos.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.FEPaisResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetTiposPaises [get]
func FEParamGetTiposPaisesHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetTiposPaises()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetActividadesHandler godoc
//
//	@Summary		Consulta Actividades
//	@Description	Esta operación permite consultar los códigos de actividades, sus descripciones y el orden (si es actividad primaria, secundaria, etc)
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wsfe.FEActividadesResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetActividades [get]
func FEParamGetActividadesHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsfe.FEParamGetActividades()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// FEParamGetCondicionIvaReceptorHandler godoc
//
//	@Summary		Consulta condiciones IVA del receptor
//	@Description	Esta operación permite consultar los identificadores de la condicion frente al IVA del receptor, su descripción y a la clase de comprobante que corresponde.
//	@Tags			Factura Electrónica
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Param			ClaseCmp	query		string	true	"Clase de Comprobate"
//	@Success		200			{object}	wsfe.CondicionIvaReceptorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/fe/FEParamGetCondicionIvaReceptor [get]
func FEParamGetCondicionIvaReceptorHandler(w http.ResponseWriter, r *http.Request) {
	claseCmp := r.URL.Query().Get("ClaseCmp")
	if claseCmp == "" {
		err := errors.New("error leyendo parámetros ClaseCmp")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	resultado, err := Wsfe.FEParamGetCondicionIvaReceptor(claseCmp)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}
