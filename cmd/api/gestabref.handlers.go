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
//	@Success		200			{object}	wgestabref.WsDummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/gestabref/Dummy [get]
func DummyGesTabRefHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wsgestabref.Dummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// ConsultarFechaUltActHandler godoc
//
//	@Summary		Obtener la Fecha de última actualización de la tabla
//	@Description	Retorna la fecha de última actualización de la tabla consultada.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.FechaUltAct
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ConsultarFechaUltAct [get]
func ConsultarFechaUltActHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	resultado, err := Wsgestabref.ConsultarFechaUltAct(idReferencia)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// ListaArancelHandler godoc
//
//	@Summary		Lista Arancel
//	@Description	Retorna tabla del tipo código / descripción / opción / vigencia.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.Opciones
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaArancel [get]
func ListaArancelHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	resultado, err := Wsgestabref.ListaArancel(idReferencia)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// ListaDescripcionHandler godoc
//
//	@Summary		Lista Descripción
//	@Description	Emite tabla del tipo código / descripción.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.Descripciones
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaDescripcion [get]
func ListaDescripcionHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaDescripcion(idReferencia)
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
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.DescripcionesCodificaciones
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaDescripcionDecodificacion [get]
func ListaDescripcionDecodificacionHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaDescripcionDecodificacion(idReferencia)
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
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.Empresas
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaEmpresas [get]
func ListaEmpresasHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaEmpresas(idReferencia)
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
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.LugaresOperativos
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaLugaresOperativos [get]
func ListaLugaresOperativosHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaLugaresOperativos(idReferencia)
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
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.PaisesAduanas
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaPaisesAduanas [get]
func ListaPaisesAduanasHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaPaisesAduanas(idReferencia)
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
//	@Success		200			{object}	wgestabref.TablasReferencia
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaTablasReferencia [get]
func ListaTablasReferenciaHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Wsgestabref.ListaTablasReferencia()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaVigenciasHandler godoc
//
//	@Summary		Lista de Vigencias
//	@Description	Emite tabla del tipo código / descripción / vigencia.
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.Vigencias
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaVigencias [get]
func ListaVigenciasHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaVigencias(idReferencia)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}

// ListaDatoComplementarioHandler godoc
//
//	@Summary		Lista Datos Complementarios
//	@Description	Lista Datos Complementarios
//	@Tags			Consulta de Tablas de Referencia
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key		header		string	true	"API Key de acceso"
//	@Param			IdReferencia	query		string	true	"ID de referencia"
//	@Success		200				{object}	wgestabref.DatosComplementarios
//	@Failure		400				{object}	dto.ErrorResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Router			/gestabref/ListaDatoComplementario [get]
func ListaDatoComplementarioHandler(w http.ResponseWriter, r *http.Request) {
	idReferencia := r.URL.Query().Get("IdReferencia")
	if len(idReferencia) <= 0 {
		err := errors.New("error leyendo parámetro IdReferencia")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	data, err := Wsgestabref.ListaDatoComplementario(idReferencia)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, data, nil)
}
