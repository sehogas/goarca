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
	var data []*wsfe.FECAEDetResponse = make([]*wsfe.FECAEDetResponse, 0)
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
			data = response.FECAESolicitarResult.FeDetResp.FECAEDetResponse
		}
	}

	// if len(data) > 0 {
	// 	resultado = data[0].Resultado
	// 	cae = data[0].CAE
	// 	caeFchVto = data[0].CAEFchVto
	// }

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposCbte() ([]*wsfe.CbteTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposCbte(&wsfe.FEParamGetTiposCbte{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.CbteTipo = make([]*wsfe.CbteTipo, 0)
	if response.FEParamGetTiposCbteResult.Events != nil {
		for _, e := range response.FEParamGetTiposCbteResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposCbteResult.Errors != nil {
		for _, e := range response.FEParamGetTiposCbteResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposCbteResult.ResultGet != nil {
			data = response.FEParamGetTiposCbteResult.ResultGet.CbteTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposConcepto() ([]*wsfe.ConceptoTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposConcepto(&wsfe.FEParamGetTiposConcepto{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.ConceptoTipo = make([]*wsfe.ConceptoTipo, 0)
	if response.FEParamGetTiposConceptoResult.Events != nil {
		for _, e := range response.FEParamGetTiposConceptoResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposConceptoResult.Errors != nil {
		for _, e := range response.FEParamGetTiposConceptoResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposConceptoResult.ResultGet != nil {
			data = response.FEParamGetTiposConceptoResult.ResultGet.ConceptoTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposDoc() ([]*wsfe.DocTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposDoc(&wsfe.FEParamGetTiposDoc{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.DocTipo = make([]*wsfe.DocTipo, 0)
	if response.FEParamGetTiposDocResult.Events != nil {
		for _, e := range response.FEParamGetTiposDocResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposDocResult.Errors != nil {
		for _, e := range response.FEParamGetTiposDocResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposDocResult.ResultGet != nil {
			data = response.FEParamGetTiposDocResult.ResultGet.DocTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposIva() ([]*wsfe.IvaTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposIva(&wsfe.FEParamGetTiposIva{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.IvaTipo = make([]*wsfe.IvaTipo, 0)
	if response.FEParamGetTiposIvaResult.Events != nil {
		for _, e := range response.FEParamGetTiposIvaResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposIvaResult.Errors != nil {
		for _, e := range response.FEParamGetTiposIvaResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposIvaResult.ResultGet != nil {
			data = response.FEParamGetTiposIvaResult.ResultGet.IvaTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposMonedas() ([]*wsfe.Moneda, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposMonedas(&wsfe.FEParamGetTiposMonedas{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.Moneda = make([]*wsfe.Moneda, 0)
	if response.FEParamGetTiposMonedasResult.Events != nil {
		for _, e := range response.FEParamGetTiposMonedasResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposMonedasResult.Errors != nil {
		for _, e := range response.FEParamGetTiposMonedasResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposMonedasResult.ResultGet != nil {
			data = response.FEParamGetTiposMonedasResult.ResultGet.Moneda
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposOpcional() ([]*wsfe.OpcionalTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposOpcional(&wsfe.FEParamGetTiposOpcional{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.OpcionalTipo = make([]*wsfe.OpcionalTipo, 0)
	if response.FEParamGetTiposOpcionalResult.Events != nil {
		for _, e := range response.FEParamGetTiposOpcionalResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposOpcionalResult.Errors != nil {
		for _, e := range response.FEParamGetTiposOpcionalResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposOpcionalResult.ResultGet != nil {
			data = response.FEParamGetTiposOpcionalResult.ResultGet.OpcionalTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetTiposTributos() ([]*wsfe.TributoTipo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposTributos(&wsfe.FEParamGetTiposTributos{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.TributoTipo = make([]*wsfe.TributoTipo, 0)
	if response.FEParamGetTiposTributosResult.Events != nil {
		for _, e := range response.FEParamGetTiposTributosResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetTiposTributosResult.Errors != nil {
		for _, e := range response.FEParamGetTiposTributosResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetTiposTributosResult.ResultGet != nil {
			data = response.FEParamGetTiposTributosResult.ResultGet.TributoTipo
		}
	}

	return data, errors.Join(errs...)
}

func (ws *Wsfe) FEParamGetPtosVenta() ([]*wsfe.PtoVenta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetPtosVenta(&wsfe.FEParamGetPtosVenta{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	var data []*wsfe.PtoVenta = make([]*wsfe.PtoVenta, 0)
	if response.FEParamGetPtosVentaResult.Events != nil {
		for _, e := range response.FEParamGetPtosVentaResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FEParamGetPtosVentaResult.Errors != nil {
		for _, e := range response.FEParamGetPtosVentaResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		if response.FEParamGetPtosVentaResult.ResultGet != nil {
			data = response.FEParamGetPtosVentaResult.ResultGet.PtoVenta
		}
	}

	return data, errors.Join(errs...)
}
