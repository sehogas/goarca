package main

import (
	"errors"
	"net/http"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

// DummyGesTabRefHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Consulta de Tablas de Referencia
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/gestabref/dummy [get]
func DummyGesTabRefHandler(w http.ResponseWriter, r *http.Request) {
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

// ConsultarFechaUltActHandler godoc
//
//	@Summary		Obtener la Fecha de última actualización de la tabla
//	@Description	Retorna la fecha de última actualización de la tabla consultada.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{object}	dto.FecUltActResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/consultar-fecha-ult-act [get]
func ConsultarFechaUltActHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	fecha, err := Wsgestabref.ConsultarFechaUltAct(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.FecUltActResponse{
		FechaUltAct: *fecha,
	}, nil)
}

// ListaArancelHandler godoc
//
//	@Summary		Lista Arancel
//	@Description	Retorna tabla del tipo código / descripción / opción / vigencia.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.Opcion
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-arancel [get]
func ListaArancelHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaArancel(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaDescripcionHandler godoc
//
//	@Summary		Lista Descripción
//	@Description	Emite tabla del tipo código / descripción.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.Descripcion
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-descripcion [get]
func ListaDescripcionHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaDescripcion(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaDescripcionDecodificacionHandler godoc
//
//	@Summary		Lista Descripción Decodificación
//	@Description	Lista Descripción Decodificación
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.DescripcionCodificacion
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-descripcion-decodificacion [get]
func ListaDescripcionDecodificacionHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaDescripcionDecodificacion(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaEmpresasHandler godoc
//
//	@Summary		Lista de Empresas
//	@Description	Emite tabla del tipo cuit / razón social.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.Empresa
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-empresas [get]
func ListaEmpresasHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaEmpresas(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaLugaresOperativosHandler godoc
//
//	@Summary		Lista de Lugares Operativos
//	@Description	Emite tabla del tipo código / descripción / vigencia / aduana / lugar operativo.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.LugarOperativo
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-lugares-operativos [get]
func ListaLugaresOperativosHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaListaLugaresOperativos(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaPaisesAduanasHandler godoc
//
//	@Summary		Lista de Paises y Aduanas
//	@Description	Emite tabla del tipo código / descripción / vigencia /país o aduana.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			argNombreTabla	query		string	true	"Nombre de la tabla"
//	@Success		200				{array}		wgestabref.PaisAduana
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-paises-aduanas [get]
func ListaPaisesAduanasHandler(w http.ResponseWriter, r *http.Request) {
	argNombreTabla := r.URL.Query().Get("argNombreTabla")
	if len(argNombreTabla) <= 2 {
		err := errors.New("error leyendo parámetros de la solicitud")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaPaisesAduanas(argNombreTabla)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaTablasReferenciaHandler godoc
//
//	@Summary		Lista de Tablas de Referencia
//	@Description	Emite tabla del tipo: Tabla de Referencia / Descripción Tabla Referencia / WebMethod (que se debe utilizar para obtener los datos de dicha tabla).
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{array}		wgestabref.TablaReferencia
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/gestabref/lista-tablas-referencia [get]
func ListaTablasReferenciaHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Wsgestabref.ListaTablasReferencia()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}
