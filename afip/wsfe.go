package afip

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/goarca/ws/wsfe"
)

type Wsfe struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
}

func NewWsfe(environment Environment, cuit int64) (*Wsfe, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSFEProduction
	} else {
		url = URLWSFETesting
	}

	return &Wsfe{
		serviceName: "wsfe",
		environment: environment,
		url:         url,
		cuit:        cuit,
	}, nil
}

func (ws *Wsfe) PrintlnAsXML(obj interface{}) {
	if ws.environment == TESTING {
		data, err := xml.MarshalIndent(obj, " ", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
	}
}

func (ws *Wsfe) FEDummy() (appServer, authServer, DbServer string, err error) {
	request := &wsfe.FEDummy{}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(client)

	response, err := service.FEDummy(request)
	if err != nil {
		return "", "", "", err
	}

	PrintlnAsXML(response)

	if response.FEDummyResult != nil {
		return response.FEDummyResult.AppServer, response.FEDummyResult.AuthServer, response.FEDummyResult.DbServer, nil
	}

	return "", "", "", nil
}

func (ws *Wsfe) FEUltimoComprobanteEmitido(ptoVta int32, cbteTipo int32) (int32, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return 0, err
	}

	request := &wsfe.FECompUltimoAutorizado{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		PtoVta:   ptoVta,
		CbteTipo: cbteTipo,
	}

	PrintlnAsXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECompUltimoAutorizado(request)
	if err != nil {
		return 0, err
	}
	PrintlnAsXML(response)

	var errs []error
	if response.FECompUltimoAutorizadoResult.Events != nil {
		for _, e := range response.FECompUltimoAutorizadoResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FECompUltimoAutorizadoResult.Errors != nil {
		for _, e := range response.FECompUltimoAutorizadoResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		cbteNro := response.FECompUltimoAutorizadoResult.CbteNro
		return cbteNro, errors.Join(errs...)
	}

	return 0, errors.Join(errs...)
}

func (ws *Wsfe) FECAESolicitar(cab *wsfe.FECabRequest, det []*wsfe.FECAEDetRequest) ([]*wsfe.FECAEDetResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wsfe.FECAESolicitar{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		FeCAEReq: &wsfe.FECAERequest{
			FeCabReq: &wsfe.FECAECabRequest{
				FECabRequest: cab,
			},
			FeDetReq: &wsfe.ArrayOfFECAEDetRequest{
				FECAEDetRequest: det,
			},
		},
	}

	PrintlnAsXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECAESolicitar(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var resultado []*wsfe.FECAEDetResponse = make([]*wsfe.FECAEDetResponse, 0)
	if response.FECAESolicitarResult.Events != nil {
		for _, e := range response.FECAESolicitarResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FECAESolicitarResult.Errors != nil {
		for _, e := range response.FECAESolicitarResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FECAESolicitarResult.FeDetResp != nil {
			//resultado := response.FECAESolicitarResult.FeCabResp.Resultado
			resultado = response.FECAESolicitarResult.FeDetResp.FECAEDetResponse
		}
	}

	return resultado, errors.Join(errs...)
}
