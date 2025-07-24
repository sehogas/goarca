package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
	"github.com/sehogas/goarca/ws/wscoem"
)

// DummyHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Comunicación de Embarque
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	wscoem.ResultadoEjecucionDummy
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//
//	@Router			/coem/Dummy [get]
func DummyCoemHandler(w http.ResponseWriter, r *http.Request) {
	resultado, err := Wscoem.Dummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// RegistrarCaratulaHandler godoc
//
//	@Summary		Registrar Carátula
//	@Description	Crea el identificador necesario para inicializar el circuito de comunicación de embarque.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		wscoem.RegistrarCaratulaRequest	true	"RegistrarCaratulaRequest"
//	@Success		200			{object}	dto.MessageResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/RegistrarCaratula [post]
func RegistrarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.RegistrarCaratulaRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	log.Println(post)

	resultado, err := Wscoem.RegistrarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// AnularCaratulaHandler godoc
//
//	@Summary		Anular Carátula
//	@Description	Método que permite eliminar una carátula. Si se encuentran COEMs en estado Presentada o Autorizada se detiene la ejecución.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		wscoem.AnularCaratulaRequest	true	"AnularCaratulaRequest"
//	@Success		200			{object}	wscoem.AnularEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/AnularCaratula [delete]
func AnularCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.AnularCaratulaRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.AnularCaratula(&wscoem.AnularCaratulaRequest{
		IdentificadorCaratula: post.IdentificadorCaratula,
	})
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// RectificarCaratulaHandler godoc
//
//	@Summary		Rectificar Carátula
//	@Description	Permite rectificar los datos de una Carátula previamente creada con el método RegistrarCaratula. Entre las restricciones, no se permitirá cargar datos idénticos a otra Carátula existente, ni modificar aquella carátula que tenga COEMs asociados.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string								true	"API Key de acceso"
//	@Param			request		body		wscoem.RectificarCaratulaRequest	true	"RectificarCaratulaRequest"
//	@Success		200			{object}	wscoem.RectificarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/RectificarCaratula [put]
func RectificarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.RectificarCaratulaRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.RectificarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// RegistrarCOEMHandler godoc
//
//	@Summary		Registrar COEM
//	@Description	Método a través del cual se registran los valores de una COEM comprendidos en información de Contenedores con Carga, Contenedores Vacíos y Mercadería suelta, asociados a un Identificador de Carátula previamente creado con el método RegistrarCaratula.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		wscoem.RegistrarCOEMRequest	true	"RegistrarCOEMRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/RegistrarCOEM [post]
func RegistrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.RegistrarCOEMRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.RegistrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarCambioBuqueHandler godoc
//
//	@Summary		Solicitar cambio de Buque
//	@Description	Método a través del cual se modificarán el Identificador y/o nombre del buque. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string								true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarCambioBuqueRequest	true	"SolicitarCambioBuqueRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarCambioBuque [put]
func SolicitarCambioBuqueHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarCambioBuqueRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.SolicitarCambioBuque(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarCambioFechasHandler godoc
//
//	@Summary		Solicitar cambio de Fechas
//	@Description	Método a través del cual se modifican las Fechas de Arribo y/o Fecha de Zarpada de la Carátula. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string								true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarCambioFechasRequest	true	"SolicitarCambioFechasRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarCambioFechas [put]
func SolicitarCambioFechasHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarCambioFechasRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.SolicitarCambioFechas(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarCambioLOTHandler godoc
//
//	@Summary		Solicitar cambio de LOT
//	@Description	Método a través del cual se modifica el Lugar Operativo de la Carátula. Deben existir COEMs presentadas o autorizadas, caso contrario aún se puede enviar la rectificación de la carátula. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string								true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarCambioLOTRequest	true	"SolicitarCambioLOTRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarCambioLOT [put]
func SolicitarCambioLOTHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarCambioLOTRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}
	resultado, err := Wscoem.SolicitarCambioLOT(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// RectificarCOEMHandler godoc
//
//	@Summary		Rectificar COEM
//	@Description	Método a través del cual se modifican los valores de una COEM. La COEM debe estar en curso o registrada. No se habilita rectificar una COEM presentada o autorizada. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		wscoem.RectificarCOEMRequest	true	"RectificarCOEMRequest"
//	@Success		200			{object}	wscoem.RectificarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/RectificarCOEM [put]
func RectificarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.RectificarCOEMRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.RectificarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// CerrarCOEMHandler godoc
//
//	@Summary		Cerrar COEM
//	@Description	Método a través del cual se cierra la carga de una COEM asociada a una Carátula, permitiendo que el operador Aduanero pueda trabajar con ella. La COEM debe estar en estado CUR.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		wscoem.CerrarCOEMRequest	true	"CerrarCOEMRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/CerrarCOEM [post]
func CerrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.CerrarCOEMRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.CerrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// AnularCOEMHandler godoc
//
//	@Summary		Anular COEM
//	@Description	Método a través del cual se anula una COEM. La COEM debe estar en estado en CURSO o REGISTRADA.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string						true	"API Key de acceso"
//	@Param			request		body		wscoem.AnularCOEMRequest	true	"AnularCOEMRequest"
//	@Success		200			{object}	wscoem.AnularEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/AnularCOEM [delete]
func AnularCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.AnularCOEMRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.AnularCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarAnulacionCOEMHandler godoc
//
//	@Summary		Solicitar Anulación COEM
//	@Description	Método a través del cual se solicita la anulación de una COEM cuando ya se encuentra en estado PRESENTADA o AUTORIZADA, caso contrario puede utilizar el método AnularCOEM.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string									true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarAnulacionCOEMRequest	true	"SolicitarAnulacionCOEMRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarAnulacionCOEM [post]
func SolicitarAnulacionCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarAnulacionCOEMRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.SolicitarAnulacionCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarNoABordoHandler godoc
//
//	@Summary		Solicitar No Abordo
//	@Description	Método a través del cual se solicita no abordar una o varias líneas de mercadería o contenedores asociados a una carátula/COEM. La COEM debe estar en estado PRESENTADA o AUTORIZADA. No debe existir solicitud de cierre de carga iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string							true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarNoABordoRequest	true	"SolicitarNoABordoRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarNoABordo [post]
func SolicitarNoABordoHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarNoABordoRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.SolicitarNoABordo(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarCierreCargaContoBultoHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Contenedores y/o Bultos
//	@Description	Método a través del cual se solicita el cierre de carga de una Carátula que contiene COMEs cuyos permisos de embarque amparan contenedores o bultos sueltos (no Granel). Todas las COEMs deben estar en estado AUTORIZADA o ANULADA. No se consideran las marcadas como NO ABORDO. No debe existir solicitud de cierre de carga contenedores o bultos sueltos iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string											true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarCierreCargaContoBultoRequest	true	"SolicitarCierreCargaContoBultoRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarCierreCargaContoBulto [post]
func SolicitarCierreCargaContoBultoHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarCierreCargaContoBultoRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	resultado, err := Wscoem.SolicitarCierreCargaContoBulto(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, resultado, nil)
}

// SolicitarCierreCargaGranelHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Granel
//	@Description	Método a través del cual se solicita el cierre de carga de una Carátula que contiene COEMs cuyos permisos de embarque contienen Mercadería a granel. Todas las COEMs deben estar en estado AUTORIZADA o ANULADA. No se consideran las líneas marcadas como NO ABORDO. Requiere detalle de todas las COMEs, Permisos de Embarque e ítems abordo. No debe existir solicitud de cierre de carga granel iniciada o aprobada.
//	@Tags			Comunicación de Embarque
//	@Accept			json
//	@Produce		json
//	@Param			x-api-key	header		string										true	"API Key de acceso"
//	@Param			request		body		wscoem.SolicitarCierreCargaGranelRequest	true	"SolicitarCierreCargaGranelRequest"
//	@Success		200			{object}	wscoem.RegistrarEmbarqueRta
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/coem/SolicitarCierreCargaGranel [post]
func SolicitarCierreCargaGranelHandler(w http.ResponseWriter, r *http.Request) {
	var post wscoem.SolicitarCierreCargaGranelRequest
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	respuesta, err := Wscoem.SolicitarCierreCargaGranel(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, respuesta, nil)
}
