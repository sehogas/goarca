package dto

import (
	"time"

	"github.com/sehogas/goarca/ws/wsfe"
)

type DummyResponse struct {
	AppServer  string `json:"AppServer"`
	AuthServer string `json:"AuthServer"`
	DbServer   string `json:"DbServer"`
}

type FecUltActResponse struct {
	FechaUltAct time.Time `json:"FechaUltAct"`
}

type CbteNroResponse struct {
	CbteNro int32 `json:"CbteNro"`
}

type FECAESolicitarRequest struct {
	Cab *wsfe.FECabRequest      `json:"Cabecera"`
	Det []*wsfe.FECAEDetRequest `json:"Detalle"`
}
