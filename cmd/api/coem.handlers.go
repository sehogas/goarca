package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sehogas/goarca/afip"
	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
	"github.com/sehogas/goarca/internal/util/validator"
)

// DummyHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Comunicación de Embarque
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//
//	@Router			/coem/dummy [get]
func DummyCoemHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wscoem.Dummy()
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

// RegistrarCaratulaHandler godoc
//
//	@Summary		Registrar Carátula
//	@Description	Crea el identificador necesario para inicializar el circuito de comunicación de embarque.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string				true	"API Key de acceso"
//	@Param			request		body		afip.CaratulaParams	true	"RegistrarCaratulaRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/registrar-caratula [post]
func RegistrarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RegistrarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RectificarCaratulaHandler godoc
//
//	@Summary		Rectificar Carátula
//	@Description	Permite rectificar los datos de una Carátula previamente creada con el método RegistrarCaratula. Entre las restricciones, no se permitirá cargar datos idénticos a otra Carátula existente, ni modificar aquella carátula que tenga COEMs asociados.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		afip.RectificarCaratulaParams	true	"RectificarCaratulaRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/rectificar-caratula [put]
func RectificarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RectificarCaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RectificarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// AnularCaratulaHandler godoc
//
//	@Summary		Anular Carátula
//	@Description	Método que permite eliminar una carátula. Si se encuentran COEMs en estado Presentada o Autorizada se detiene la ejecución.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string								true	"API Key de acceso"
//	@Param			request		body		afip.IdentificadorCaraturaParams	true	"AnularCaratulaRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/anular-caratula [delete]
func AnularCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCaraturaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.AnularCaratula(&afip.IdentificadorCaraturaParams{
		IdentificadorCaratula: post.IdentificadorCaratula,
	})
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioBuqueHandler godoc
//
//	@Summary		Solicitar cambio de Buque
//	@Description	Método a través del cual se modificarán el Identificador y/o nombre del buque. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string					true	"API Key de acceso"
//	@Param			request		body		afip.CambioBuqueParams	true	"SolicitarCambioBuqueRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-cambio-buque [put]
func SolicitarCambioBuqueHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioBuqueParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioBuque(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioFechasHandler godoc
//
//	@Summary		Solicitar cambio de Fechas
//	@Description	Método a través del cual se modifican las Fechas de Arribo y/o Fecha de Zarpada de la Carátula. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string					true	"API Key de acceso"
//	@Param			request		body		afip.CambioFechasParams	true	"CambioFechasParamsRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-cambio-fechas [put]
func SolicitarCambioFechasHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioFechasParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioFechas(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioLOTHandler godoc
//
//	@Summary		Solicitar cambio de LOT
//	@Description	Método a través del cual se modifica el Lugar Operativo de la Carátula. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string					true	"API Key de acceso"
//	@Param			request		body		afip.CambioLOTParams	true	"CambioLOTParamsRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-cambio-lot [put]
func SolicitarCambioLOTHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioLOTParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioLOT(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RegistrarCOEMHandler godoc
//
//	@Summary		Registrar COEM
//	@Description	Método a través del cual se registran los valores de una COEM comprendidos en información de Contenedores con Carga, Contenedores Vacíos y Mercadería suelta, asociados a un Identificador de Carátula previamente creado con el método RegistrarCaratula.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		afip.RegistrarCOEMParams	true	"RegistrarCOEMRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/registrar-coem [post]
func RegistrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RegistrarCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RegistrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RectificarCOEMHandler godoc
//
//	@Summary		Rectificar COEM
//	@Description	Método a través del cual se modifican los valores de una COEM. La COEM debe estar en curso o registrada. No se habilita rectificar una COEM presentada o autorizada. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		afip.RectificarCOEMParams	true	"RectificarCOEMRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/rectificar-coem [put]
func RectificarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RectificarCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RectificarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// CerrarCOEMHandler godoc
//
//	@Summary		Cerrar COEM
//	@Description	Método a través del cual se cierra la carga de una COEM asociada a una Carátula, permitiendo que el operador Aduanero pueda trabajar con ella. La COEM debe estar en estado CUR.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		afip.IdentificadorCOEMParams	true	"CerrarCOEMRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/cerrar-coem [post]
func CerrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.CerrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// AnularCOEMHandler godoc
//
//	@Summary		Anular COEM
//	@Description	Método a través del cual se anula una COEM. La COEM debe estar en estado en CURSO o REGISTRADA.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		afip.IdentificadorCOEMParams	true	"AnularCOEMRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/anular-coem [delete]
func AnularCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.AnularCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarAnulacionCOEMHandler godoc
//
//	@Summary		Solicitar Anulación COEM
//	@Description	Método a través del cual se solicita la anulación de una COEM cuando ya se encuentra en estado PRESENTADA o AUTORIZADA, caso contrario puede utilizar el método AnularCOEM.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		afip.IdentificadorCOEMParams	true	"SolicitarAnulacionCOEMRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-anulacion-coem [post]
func SolicitarAnulacionCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarAnulacionCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarNoABordoHandler godoc
//
//	@Summary		Solicitar No Abordo
//	@Description	Método a través del cual se solicita no abordar una o varias líneas de mercadería o contenedores asociados a una carátula/COEM. La COEM debe estar en estado PRESENTADA o AUTORIZADA. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		afip.IdentificadorCOEMParams	true	"SolicitarNoABordoRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-no-abordo [post]
func SolicitarNoABordoHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarNoABordoParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarNoABordo(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCierreCargaContoBultoHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Contenedores y/o Bultos
//	@Description	Método a través del cual se solicita el cierre de carga de una Carátula que contiene COMEs cuyos permisos de embarque amparan contenedores o bultos sueltos (no Granel). Todas las COEMs deben estar en estado AUTORIZADA o ANULADA. No se consideran las marcadas como NO ABORDO. No debe existir solicitud de cierre de carga contenedores o bultos sueltos iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string										true	"API Key de acceso"
//	@Param			request		body		afip.SolicitarCierreCargaContoBultoParams	true	"SolicitarCierreCargaContoBultoRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-cierre-carga-conto-bulto [post]
func SolicitarCierreCargaContoBultoHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarCierreCargaContoBultoParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCierreCargaContoBulto(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCierreCargaGranelHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Granel
//	@Description	Método a través del cual se solicita el cierre de carga de una Carátula que contiene COEMs cuyos permisos de embarque contienen Mercadería a granel. Todas las COEMs deben estar en estado AUTORIZADA o ANULADA. No se consideran las líneas marcadas como NO ABORDO. Requiere detalle de todas las COMEs, Permisos de Embarque e ítems abordo. No debe existir solicitud de cierre de carga granel iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string									true	"API Key de acceso"
//	@Param			request		body		afip.SolicitarCierreCargaGranelParams	true	"SolicitarCierreCargaGranel"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/solicitar-cierre-carga-granel [post]
func SolicitarCierreCargaGranelHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarCierreCargaGranelParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCierreCargaGranel(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}
