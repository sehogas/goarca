package main

import (
	"net/http"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

// InfoHandler godoc
//
//	@Summary		Muesta información de la API
//	@Description	Muesta información de la API
//	@Tags			API
//	@Produce		json
//	@Param			x-api-key	header		string	true	"API Key de acceso"
//	@Success		200			{object}	dto.DummyResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Router			/info [get]
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	util.HttpResponseJSON(w, http.StatusOK, &dto.InfoResponse{
		Version: Version,
	}, nil)
}
