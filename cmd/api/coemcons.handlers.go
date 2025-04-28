package main

import (
	"errors"
	"net/http"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

// DummyCoemconsHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Consultas de Comunicación de Embarque
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coemcons/dummy [get]
func DummyCoemconsHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wscoemcons.Dummy()
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

// ObtenerConsultaEstadosCOEMHandler godoc
//
//	@Summary		Obtener Consulta Estados COEM
//	@Description	Obtiene una lista de comunicaciones de embarque comprendidas en la cabecera informada.
//	@Tags			Consultas de Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			identificador	query		string	true	"Identificador de la caratula"
//	@Success		200				{array}		wscoemcons.ConsultaEstadoCOEM
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/coemcons/obtener-consulta-estados-coem [get]
func ObtenerConsultaEstadosCOEMHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaEstadosCOEM(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ObtenerConsultaNoAbordoHandler godoc
//
//	@Summary		Obtener Consulta No Abordo
//	@Description	Obtiene una lista de aquellas comunicaciones de embarque relacionadas con la carátula de referencia para las cuales se ha definido el no abordaje parcial o total de su contenido, el cual puede tratarse de bultos sueltos, mercadería a granel, contenedores de carga o cont¢nedores vacíos.
//	@Tags			Consultas de Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			identificador	query		string	true	"Identificador de la caratula"
//	@Success		200				{array}		wscoemcons.ConsultaNoAbordo
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/coemcons/obtener-consulta-no-abordo [get]
func ObtenerConsultaNoAbordoHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaNoAbordo(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ObtenerConsultaSolicitudesHandler godoc
//
//	@Summary		Obtener Consulta de Solicitudes
//	@Description	Obtiene la lista de solicitudes de comunicaciones asociadas a la cabecera informada. Número y tipo de solicitud, estado y fecha.
//	@Tags			Consultas de Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			identificador	query		string	true	"Identificador de la caratula"
//	@Success		200				{array}		wscoemcons.ConsultaSolicitudes
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/coemcons/obtener-consulta-solicitudes [get]
func ObtenerConsultaSolicitudesHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaSolicitudes(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}
